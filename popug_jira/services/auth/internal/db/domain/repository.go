package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user *User) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, userName string) (*User, error)
}
