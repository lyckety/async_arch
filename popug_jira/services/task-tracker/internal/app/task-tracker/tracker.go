package producer

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	tasksEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/events/tasks"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/interceptors"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	pbV1Task "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/task/v1"
	pbV1TaskEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/taskevents/v1"
	pbV1Tasks "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/tasktracker/v1"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskTrackerService struct {
	pbV1Tasks.UnimplementedTaskTrackerServiceServer

	dbIns domain.Repository

	producerEvents *tasksEvents.TaskCUDEventSender
}

func New(
	db domain.Repository,
	eventSender *tasksEvents.TaskCUDEventSender,
) *TaskTrackerService {
	return &TaskTrackerService{
		dbIns:          db,
		producerEvents: eventSender,
	}
}

func (s *TaskTrackerService) CreateTask(
	ctx context.Context,
	req *pbV1Tasks.CreateTaskRequest,
) (*pbV1Tasks.CreateTaskResponse, error) {
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

	taskID, err := s.dbIns.CreateTask(ctx, newTaskDB)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error create task %v: %s", newTaskDB, err.Error(),
		)
	}

	rpcTaskInfo := &pbV1Task.Task{
		Description:  req.GetDescription(),
		AssignedUser: assignedUser.ID.String(),
		Status:       pbV1Task.TaskStatus_TASK_STATUS_OPENED,
	}

	eventsMsg := &pbV1TaskEvents.TaskEvent{
		EventType: pbV1TaskEvents.TaskCUDEventType_TASK_CUD_EVENT_TYPE_CREATED,
		Task: &pbV1Task.TaskWithID{
			Id:   taskID.String(),
			Task: rpcTaskInfo,
		},
		Timestamp: newTaskDB.CreatedAt.Unix(), // TODO: FIX ME !!!! надо проверить, что везде он проставляется со значением из БД
	}

	go func() {
		if err := s.producerEvents.Send(ctx, eventsMsg); err != nil {
			log.Errorf("failed send event for create task %v: %s", rpcTaskInfo, err.Error())
		}

		log.Debugf("success sent cud event (created task): %v!", eventsMsg)
	}()

	return &pbV1Tasks.CreateTaskResponse{
		Id: taskID.String(),
	}, nil
}

func (s *TaskTrackerService) RandomReassignOpenedTasks(
	ctx context.Context,
	req *pbV1Tasks.RandomReassignOpenedTasksRequest,
) (*pbV1Tasks.RandomReassignOpenedTasksResponse, error) {
	tasksToUsers, err := s.dbIns.RandomlyUpdateAssignedOpenedTasks(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"error random assigned users to all opened tasks: %s", err.Error(),
		)
	}

	rpcResponse := make([]*pbV1Tasks.TaskIDAndAssignedUserID, len(tasksToUsers))
	cudEvents := make([]*pbV1TaskEvents.TaskEvent, len(tasksToUsers))

	cnt := 0
	for reassignedTask := range tasksToUsers {
		rpcResponse[cnt] = &pbV1Tasks.TaskIDAndAssignedUserID{
			TaskId:         reassignedTask.ID.String(),
			AssignedUserId: reassignedTask.UserID.String(),
		}

		cudEvents[cnt] = &pbV1TaskEvents.TaskEvent{
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
		if err := s.producerEvents.Send(ctx, cudEvents...); err != nil {
			log.Errorf("failed send event for random reassigned tasks: %s", err.Error())
		}

		log.Debugf("success sent cud event (random tasks reassigned)!")
	}()

	return &pbV1Tasks.RandomReassignOpenedTasksResponse{
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

	task, err := s.dbIns.TaskCompleteByUser(ctx, userUUID, taskUUID)
	if err != nil {
		log.Errorf("TaskComplete(...): task id must be uuid: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"TaskComplete(...): task id must be uuid: %s", err.Error(),
		)
	}

	cudEvent := &pbV1TaskEvents.TaskEvent{
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
		if err := s.producerEvents.Send(ctx, cudEvent); err != nil {
			log.Errorf("failed send event for complete task %q: %s", cudEvent, err.Error())
		}

		log.Debugf("success sent cud event (task copleted): %v", cudEvent)
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

func (s *TaskTrackerService) convertRPCTaskStatusToDB(rpcTaskStatus pbV1Task.TaskStatus) domain.TaskStatus {
	switch rpcTaskStatus {
	case pbV1Task.TaskStatus_TASK_STATUS_OPENED:
		return domain.TaskOpened
	case pbV1Task.TaskStatus_TASK_STATUS_COMPLETED:
		return domain.TaskCompleted
	case pbV1Task.TaskStatus_TASK_STATUS_UNSPECIFIED:
		fallthrough
	default:
		return domain.TaskUnknowStatus
	}
}

func (s *TaskTrackerService) getRandomUser(users []*domain.User) *domain.User {
	index := rand.Intn(len(users))

	return users[index]
}
