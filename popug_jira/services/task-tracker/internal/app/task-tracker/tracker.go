package producer

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	tasksEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/events/tasks"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/interceptors"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	pbV1Task "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/api/grpc/task/v1"
	pbV1Tasks "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/api/grpc/tasktracker/v1"
	pbV1TaskEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/taskevents/v1"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskTrackerService struct {
	pbV1Tasks.UnimplementedTaskTrackerServiceServer

	dbIns domain.Repository

	producerCUDEvents      *tasksEvents.TaskCUDEventSender
	producerBusinessEvents *tasksEvents.TaskCUDEventSender
}

func New(
	db domain.Repository,
	cudSender *tasksEvents.TaskCUDEventSender,
	businessSender *tasksEvents.TaskCUDEventSender,
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
	newTaskDB := &domain.Task{
		Description: req.GetDescription(),
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

	eventsMsg := &pbV1TaskEvents.TaskEvent{
		EventType: pbV1TaskEvents.TaskCUDEventType_TASK_CUD_EVENT_TYPE_CREATED,
		Task: &pbV1Task.TaskWithID{
			Id:   string(createdTask.PublicID.String()),
			Task: rpcTaskInfo,
		},
		Timestamp: newTaskDB.CreatedAt.Unix(),
	}

	go func() {
		if err := s.producerCUDEvents.Send(ctx, eventsMsg); err != nil {
			log.Errorf("failed send cud event for create task %v: %s", rpcTaskInfo, err.Error())

			return
		}

		log.Debugf("success sent cud event (created task): %v!", eventsMsg)
	}()

	go func() {
		if err := s.producerBusinessEvents.Send(ctx, eventsMsg); err != nil {
			log.Errorf("failed send business event for create task %v: %s", rpcTaskInfo, err.Error())

			return
		}

		log.Debugf("success sent business event (created task): %v!", eventsMsg)
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
	events := make([]*pbV1TaskEvents.TaskEvent, len(tasksToUsers))

	cnt := 0
	for reassignedTask := range tasksToUsers {
		rpcResponse[cnt] = &pbV1Tasks.TaskIDAndAssignedUserID{
			TaskId:         reassignedTask.ID.String(),
			AssignedUserId: reassignedTask.UserID.String(),
		}

		events[cnt] = &pbV1TaskEvents.TaskEvent{
			EventType: pbV1TaskEvents.TaskCUDEventType_TASK_CUD_EVENT_TYPE_ASSIGNED,
			Task: &pbV1Task.TaskWithID{
				Id: reassignedTask.ID.String(),
				Task: &pbV1Task.Task{
					Description:  reassignedTask.Description,
					AssignedUser: reassignedTask.UserID.String(),
					Status:       pbV1Task.TaskStatus_TASK_STATUS_OPENED,
				},
			},
			Timestamp: reassignedTask.UpdatedAt.Unix(),
		}

		cnt++
	}

	go func() {
		if err := s.producerCUDEvents.Send(ctx, events...); err != nil {
			log.Errorf("failed send cud events for random reassigned tasks: %s", err.Error())

			return
		}

		log.Debugf("success sent cud events (random tasks reassigned)!")
	}()

	go func() {
		if err := s.producerBusinessEvents.Send(ctx, events...); err != nil {
			log.Errorf("failed send business events for random reassigned tasks: %s", err.Error())

			return
		}

		log.Debugf("success sent business events (random tasks reassigned)!")
	}()

	return &pbV1Tasks.TasksShuffleReasignResponse{
			TaskToAssignedUser: rpcResponse,
		},
		nil
}

func (s *TaskTrackerService) TaskComplete(
	ctx context.Context,
	req *pbV1Tasks.TaskCompleteRequest,
) (*pbV1Tasks.TaskCompleteResponse, error) {
	userID := ctx.Value(interceptors.ContextKeyUserID.String()).(string)

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

	task, err := s.dbIns.TaskCompleteByUser(ctx, userUUID, taskUUID)
	if err != nil {
		log.Errorf("TaskComplete(...): task id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"TaskComplete(...): task id must be uuid: %s", err.Error(),
		)
	}

	event := &pbV1TaskEvents.TaskEvent{
		EventType: pbV1TaskEvents.TaskCUDEventType_TASK_CUD_EVENT_TYPE_COMPLETED,
		Task: &pbV1Task.TaskWithID{
			Id: task.ID.String(),
			Task: &pbV1Task.Task{
				Description:  task.Description,
				AssignedUser: task.UserID.String(),
				Status:       pbV1Task.TaskStatus_TASK_STATUS_COMPLETED,
			},
		},
		Timestamp: task.UpdatedAt.Unix(),
	}

	go func() {
		if err := s.producerCUDEvents.Send(ctx, event); err != nil {
			log.Errorf("failed send cud event for completed task %q: %s", event, err.Error())

			return
		}

		log.Debugf("success sent cud event (task completed): %v", event)
	}()

	go func() {
		if err := s.producerBusinessEvents.Send(ctx, event); err != nil {
			log.Errorf("failed send business event for random reassigned tasks: %s", err.Error())

			return
		}

		log.Debugf("success sent business event (random tasks completed)!")
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
