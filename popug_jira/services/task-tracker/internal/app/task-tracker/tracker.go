package producer

// TODO: name package

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	tasksEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/events/tasks"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/helpers"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/interceptors"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	pbV1Task "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/api/grpc/task/v1"
	pbV1Tasks "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/api/grpc/tasktracker/v1"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskTrackerService struct {
	pbV1Tasks.UnimplementedTaskTrackerServiceServer

	dbIns domain.Repository

	producerCUDEvents      *tasksEvents.TaskEventSender
	producerBusinessEvents *tasksEvents.TaskEventSender
}

func New(
	db domain.Repository,
	cudSender *tasksEvents.TaskEventSender,
	businessSender *tasksEvents.TaskEventSender,
) *TaskTrackerService {
	return &TaskTrackerService{
		dbIns:                  db,
		producerCUDEvents:      cudSender,
		producerBusinessEvents: businessSender,
	}
}

func (s *TaskTrackerService) TaskCreate(
	ctx context.Context,
	req *pbV1Tasks.TaskCreateRequest,
) (*pbV1Tasks.TaskCreateResponse, error) {
	jiraID, title := helpers.ParseTaskDescription(req.GetDescription())

	if err := helpers.ValidateJiraIDAndTitle(jiraID, title); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"error create task %v: validate jira id and title: %s: valid [<jira_id>]: <some_title>", req, err.Error(),
		)
	}

	newTaskDB := &domain.Task{
		Description: title,
		JiraId:      jiraID,
		Status:      domain.TaskOpened,
	}

	workerUsers, err := s.dbIns.GetUsersByRole(ctx, domain.WorkerRole)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error create task %v: %s", newTaskDB, err.Error(),
		)
	}

	if len(workerUsers) == 0 {
		return nil, status.Errorf(
			codes.Internal,
			"error create task %v: not found users with role %q: %s", newTaskDB, domain.WorkerRole, err,
		)
	}

	assignedUser := s.getRandomUser(workerUsers)

	newTaskDB.UserID = assignedUser.ID

	createdTask, err := s.dbIns.CreateTask(ctx, newTaskDB)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error create task %v: %s", newTaskDB, err.Error(),
		)
	}

	rpcTaskInfo := &pbV1Task.Task{
		Description:  req.GetDescription(),
		AssignedUser: string(assignedUser.ID.String()),
		Status:       pbV1Task.TaskStatus_TASK_STATUS_OPENED,
	}

	go func() {
		// eventCUDMsgV1 := tasksEvents.NewTaskCreatedV1(
		// 	createdTask.PublicID,
		// 	createdTask.UserID,
		// 	createdTask.Description,
		// 	createdTask.CreatedAt.Unix(),
		// )

		eventCUDMsgV2 := tasksEvents.NewTaskCreatedV2(
			createdTask.PublicID,
			createdTask.UserID,
			title,
			jiraID,
			createdTask.CreatedAt.Unix(),
		)

		if err := s.producerCUDEvents.Send(context.Background(), eventCUDMsgV2); err != nil {
			log.Errorf("failed send cud event for create task %v v2: %s", rpcTaskInfo, err.Error())

			return
		}

		log.Debugf("success sent cud event (created task v2): %v!", eventCUDMsgV2)
	}()

	go func() {
		eventBEMsg := tasksEvents.NewTaskAssignedV1(
			createdTask.PublicID,
			createdTask.UserID,
			createdTask.CreatedAt.Unix(),
		)

		if err := s.producerBusinessEvents.Send(context.Background(), eventBEMsg); err != nil {
			log.Errorf("failed send business event for create task %v: %s", rpcTaskInfo, err.Error())

			return
		}

		log.Debugf("success sent business event (assigned task): %v!", eventBEMsg)
	}()

	return &pbV1Tasks.TaskCreateResponse{
		Id: string(createdTask.PublicID.String()),
	}, nil
}

func (s *TaskTrackerService) TasksShuffleReasign(
	ctx context.Context,
	req *pbV1Tasks.TasksShuffleReasignRequest,
) (*pbV1Tasks.TasksShuffleReasignResponse, error) {
	tasksToUsers, err := s.dbIns.RandomlyUpdateAssignedOpenedTasks(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error random assigned users to all opened tasks: %s", err.Error(),
		)
	}

	rpcResponse := make([]*pbV1Tasks.TaskIDAndAssignedUserID, len(tasksToUsers))
	events := make([]interface{}, len(tasksToUsers))

	cnt := 0
	for _, reassignedTask := range tasksToUsers {
		rpcResponse[cnt] = &pbV1Tasks.TaskIDAndAssignedUserID{
			TaskId:         reassignedTask.ID.String(),
			AssignedUserId: reassignedTask.UserID.String(),
		}

		events[cnt] = tasksEvents.NewTaskAssignedV1(
			reassignedTask.PublicID,
			reassignedTask.UserID,
			reassignedTask.UpdatedAt.Unix(),
		)

		cnt++
	}

	if len(events) > 0 {
		go func() {
			if err := s.producerBusinessEvents.Send(context.Background(), events...); err != nil {
				log.Errorf("failed send business events for random reassigned tasks: %s", err.Error())

				return
			}

			log.Debugf("success sent business events (random tasks reassigned)!")
		}()
	}

	return &pbV1Tasks.TasksShuffleReasignResponse{
			TaskToAssignedUser: rpcResponse,
		},
		nil
}

func (s *TaskTrackerService) TaskComplete(
	ctx context.Context,
	req *pbV1Tasks.TaskCompleteRequest,
) (*pbV1Tasks.TaskCompleteResponse, error) {
	userID := ctx.Value(interceptors.ContextKeyUserID).(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Errorf("TaskComplete(...): user id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"TaskComplete(...): user id must be uuid: %s", err.Error(),
		)
	}

	taskUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Errorf("task id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"TaskComplete(...): task id must be uuid: %s", err.Error(),
		)
	}

	task, err := s.dbIns.TaskCompleteByUser(context.Background(), userUUID, taskUUID)
	if err != nil {
		log.Errorf("TaskComplete(...): task id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"TaskComplete(...): task id must be uuid: %s", err.Error(),
		)
	}

	eventsMsg := tasksEvents.NewTaskCompletedV1(
		task.PublicID,
		task.UserID,
		task.UpdatedAt.Unix(),
	)

	go func() {
		if err := s.producerBusinessEvents.Send(context.Background(), eventsMsg); err != nil {
			log.Errorf("failed send business event completed task: %s", err.Error())

			return
		}

		log.Debugf("success sent business event (tasks completed)!")
	}()

	return &pbV1Tasks.TaskCompleteResponse{}, nil
}

func (s *TaskTrackerService) GetListOpenedTasksForMe(
	ctx context.Context,
	req *pbV1Tasks.GetListOpenedTasksForMeRequest,
) (*pbV1Tasks.GetListOpenedTasksForMeResponse, error) {
	userID := ctx.Value(interceptors.ContextKeyUserID).(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Errorf("GetListOpenedTasksForMe(...): user id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"GetListOpenedTasksForMe(...): user id must be uuid: %s", err.Error(),
		)
	}

	tasksDB, err := s.dbIns.GetAllTasksByUserAndStatus(ctx, userUUID, domain.TaskOpened)
	if err != nil {
		log.Errorf(
			"s.dbIns.GetAllTasksByUserAndStatus(ctx, %s, opened): %s",
			userUUID.String(),
			err.Error(),
		)

		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.GetAllTasksByUserAndStatus(ctx, %s, opened): %s",
			userUUID.String(),
			err.Error(),
		)
	}

	rpcTasks := make([]*pbV1Task.TaskWithID, len(tasksDB))

	for i, task := range tasksDB {
		rpcTask := &pbV1Task.TaskWithID{
			Id: task.ID.String(),
			Task: &pbV1Task.Task{
				Description:  task.Description,
				AssignedUser: task.UserID.String(),
				Status:       pbV1Task.TaskStatus_TASK_STATUS_OPENED,
			},
		}

		rpcTasks[i] = rpcTask
	}

	return &pbV1Tasks.GetListOpenedTasksForMeResponse{
		Task: rpcTasks,
	}, nil
}

func (s *TaskTrackerService) getRandomUser(users []*domain.User) *domain.User {
	index := rand.Intn(len(users))

	return users[index]
}
