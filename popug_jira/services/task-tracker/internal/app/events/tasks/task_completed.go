package tasks

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Completed "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/completed/v1"
)

type TaskCompletedV1 struct {
	// headers
	EventTime int64

	// data
	PublicID uuid.UUID
	WorkerID uuid.UUID
}

func NewTaskCompletedV1(pulicID, workerID uuid.UUID, time int64) *TaskCompletedV1 {
	return &TaskCompletedV1{
		EventTime: time,
		PublicID:  pulicID,
		WorkerID:  workerID,
	}
}

func (t TaskCompletedV1) name() string {
	return "TaskCompleted"
}

func (t TaskCompletedV1) version() uint {
	return 1
}

func (t *TaskCompletedV1) validate() error {
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

func (t *TaskCompletedV1) ToPB() (*pbV1Completed.Event, error) {
	if err := t.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", t.name(), t.version(), err)
	}

	return &pbV1Completed.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", t.version()),
			EventName:     t.name(),
			EventTime:     t.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Completed.Data{
			PublicId:       t.PublicID.String(),
			WorkerPublicId: t.WorkerID.String(),
		},
	}, nil
}
