package events

import (
	"fmt"

	"github.com/google/uuid"
	pbV1Header "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/eventheaders/header/v1"
	pbV1Created "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
)

type CreatedUserV1 struct {
	// headers
	EventTime int64

	// data
	PublicID  uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Role      string
	Email     string
}

func (u *CreatedUserV1) name() string {
	return "CreatedUser"
}

func (u *CreatedUserV1) version() int {
	return 1
}

func (u *CreatedUserV1) validate() error {
	if u.EventTime == 0 {
		return ErrNotSetEventTime
	}

	if u.PublicID == uuid.Nil {
		return ErrNotSetPublicID
	}

	if u.Username == "" {
		return ErrNotSetUserName
	}

	if u.Email == "" {
		return ErrNotSetEmail
	}

	if u.Role == "" {
		return ErrNotSetUserRole
	}

	return nil
}

func (u *CreatedUserV1) toPB() (*pbV1Created.Event, error) {
	if err := u.validate(); err != nil {
		return nil, fmt.Errorf("validate event %s.%d: %w", u.name(), u.version(), err)
	}

	userRole, err := convertStringRoleToRPCRoleV1(u.Role)
	if err != nil {
		return nil, fmt.Errorf("convertStringRoleToRPCRoleV1(%s): %w", u.Role, err)
	}

	return &pbV1Created.Event{
		Header: &pbV1Header.Header{
			EventId:       uuid.New().String(),
			EventVersion:  fmt.Sprintf("%d", u.version()),
			EventName:     u.name(),
			EventTime:     u.EventTime,
			EventProducer: producerName,
		},
		Data: &pbV1Created.Data{
			PublicId:  u.PublicID.String(),
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Role:      userRole,
		},
	}, nil
}
