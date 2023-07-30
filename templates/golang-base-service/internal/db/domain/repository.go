package domain

import (
	"context"
)

type Repository interface {
	CreateOrUpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userName string) error
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByName(ctx context.Context, userName string) (*User, error)
}
