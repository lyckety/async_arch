package tasks

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Paymented "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/paymented/v1"
)

type TxPaymentedV1 struct {
	// headers
	EventTime int64

	// data
	PublicID       uuid.UUID
	WorkerPublicID uuid.UUID
	Cost           uint32
}

func NewTxPaymentedV1(pulicID, workerID uuid.UUID, cost uint32, time int64) *TxDebitedV1 {
	return &TxDebitedV1{
		EventTime:      time,
		PublicID:       pulicID,
		WorkerPublicID: workerID,
		Cost:           cost,
	}
}

func (t TxPaymentedV1) name() string {
	return "TransactionPaymented"
}

func (t TxPaymentedV1) version() uint {
	return 1
}

func (t *TxPaymentedV1) validate() error {
	if t.EventTime == 0 {
		return ErrNotSetEventTime
	}

	if t.PublicID == uuid.Nil {
		return ErrNotSetPublicID
	}

	if t.WorkerPublicID == uuid.Nil {
		return ErrNotSetWorkerID
	}

	if t.Cost == 0 {
		return ErrNotSetCost
	}

	return nil
}

func (t *TxPaymentedV1) ToPB() (*pbV1Paymented.Event, error) {
	if err := t.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", t.name(), t.version(), err)
	}

	return &pbV1Paymented.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", t.version()),
			EventName:     t.name(),
			EventTime:     t.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Paymented.Data{
			PublicId:       t.PublicID.String(),
			WorkerPublicId: t.WorkerPublicID.String(),
			Cost:           t.Cost,
		},
	}, nil
}
