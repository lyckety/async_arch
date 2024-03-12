package tasks

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Assigned "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/assigned/v1"
)

type TaskAssignedV1 struct {
	// headers
	EventTime int64

	// data
	PublicID uuid.UUID
	WorkerID uuid.UUID
}

func NewTaskAssignedV1(pulicID, workerID uuid.UUID, time int64) *TaskAssignedV1 {
	return &TaskAssignedV1{
		EventTime: time,
		PublicID:  pulicID,
		WorkerID:  workerID,
	}
}

func (t TaskAssignedV1) name() string {
	return "TaskAssigned"
}

func (t TaskAssignedV1) version() uint {
	return 1
}

func (t *TaskAssignedV1) validate() error {
	if t.EventTime == 0 {
		return ErrNotSetEventTime
	}

	if t.PublicID == uuid.Nil {
		return ErrNotSetPublicID
	}

	if t.WorkerID == uuid.Nil {
		return ErrNotSetWorkerID
	}

	return nil
}

func (t *TaskAssignedV1) ToPB() (*pbV1Assigned.Event, error) {
	if err := t.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", t.name(), t.version(), err)
	}

	return &pbV1Assigned.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", t.version()),
			EventName:     t.name(),
			EventTime:     t.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Assigned.Data{
			PublicId:       t.PublicID.String(),
			WorkerPublicId: t.WorkerID.String(),
		},
	}, nil
}
