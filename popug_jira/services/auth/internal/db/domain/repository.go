package domain

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUsername(ctx context.Context, userName string) (*User, error)
}
