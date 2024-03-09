package tasks

import "errors"

var (
	ErrNotSetUserName          = errors.New("must be set username")
	ErrNotSetPublicID          = errors.New("must be set public id")
	ErrNotSetWorkerID          = errors.New("must be set worker id")
	ErrNotSetEmail             = errors.New("must be set user email")
	ErrNotSetEventProducerName = errors.New("must be set event producer name")
	ErrNotSetEventTime         = errors.New("must be set event time")
	ErrNotSetUserRole          = errors.New("must be set user role")
	ErrUnknownUserRole         = errors.New("unknown user role")
	ErrNotSetDescription       = errors.New("must be set description")
)
