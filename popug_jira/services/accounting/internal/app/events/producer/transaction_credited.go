package tasks

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Credited "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/credited/v1"
)

type TxCreditedV1 struct {
	// headers
	EventTime int64

	// data
	PublicID       uuid.UUID
	WorkerPublicID uuid.UUID
	TaskPublicID   uuid.UUID
	Cost           uint32
}

func NewTxCreditedV1(pulicID, workerID, taskID uuid.UUID, cost uint32, time int64) *TxCreditedV1 {
	return &TxCreditedV1{
		EventTime:      time,
		PublicID:       pulicID,
		WorkerPublicID: workerID,
		TaskPublicID:   taskID,
		Cost:           cost,
	}
}

func (t TxCreditedV1) name() string {
	return "TransactionCredited"
}

func (t TxCreditedV1) version() uint {
	return 1
}

func (t *TxCreditedV1) validate() error {
	if t.EventTime == 0 {
		return ErrNotSetEventTime
	}

	if t.PublicID == uuid.Nil {
		return ErrNotSetPublicID
	}

	if t.WorkerPublicID == uuid.Nil {
		return ErrNotSetWorkerID
	}

	if t.TaskPublicID == uuid.Nil {
		return ErrNotSetTaskID
	}

	if t.Cost == 0 {
		return ErrNotSetCost
	}

	return nil
}

func (t *TxCreditedV1) ToPB() (*pbV1Credited.Event, error) {
	if err := t.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", t.name(), t.version(), err)
	}

	return &pbV1Credited.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", t.version()),
			EventName:     t.name(),
			EventTime:     t.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Credited.Data{
			PublicId:       t.PublicID.String(),
			WorkerPublicId: t.WorkerPublicID.String(),
			TaskPublicId:   t.TaskPublicID.String(),
			Cost:           t.Cost,
		},
	}, nil
}
