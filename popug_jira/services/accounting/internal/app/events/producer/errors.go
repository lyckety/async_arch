package tasks

import "errors"

var (
	ErrNotSetPublicID          = errors.New("must be set public id")
	ErrNotSetWorkerID          = errors.New("must be set worker id")
	ErrNotSetTaskID            = errors.New("must be set task id")
	ErrNotSetCost              = errors.New("must be set cost")
	ErrNotSetEventProducerName = errors.New("must be set event producer name")
	ErrNotSetEventTime         = errors.New("must be set event time")
)
