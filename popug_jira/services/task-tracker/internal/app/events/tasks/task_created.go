package tasks

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Created "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v1"
)

type TaskCreatedV1 struct {
	// headers
	EventTime int64

	// data
	PublicID    uuid.UUID
	WorkerID    uuid.UUID
	Description string
}

func NewTaskCreatedV1(pulicID, workerID uuid.UUID, description string, time int64) *TaskCreatedV1 {
	return &TaskCreatedV1{
		EventTime:   time,
		PublicID:    pulicID,
		WorkerID:    workerID,
		Description: description,
	}
}

func (t TaskCreatedV1) name() string {
	return "TaskCreated"
}

func (t TaskCreatedV1) version() uint {
	return 1
}

func (t *TaskCreatedV1) validate() error {
	if t.EventTime == 0 {
		return ErrNotSetEventTime
	}

	if t.PublicID == uuid.Nil {
		return ErrNotSetPublicID
	}

	if t.WorkerID == uuid.Nil {
		return ErrNotSetPublicID
	}

	return nil
}

func (t *TaskCreatedV1) ToPB() (*pbV1Created.Event, error) {
	if err := t.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", t.name(), t.version(), err)
	}

	return &pbV1Created.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", t.version()),
			EventName:     t.name(),
			EventTime:     t.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Created.Data{
			PublicId:       t.PublicID.String(),
			WorkerPublicId: t.WorkerID.String(),
			Description:    t.Description,
		},
	}, nil
}
