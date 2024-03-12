package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	TaskRepository
	UserRepository
}

type TaskRepository interface {
	CreateTask(context.Context, *Task) (*Task, error)
	TaskCompleteByUser(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*Task, error)
	RandomlyUpdateAssignedOpenedTasks(ctx context.Context) ([]Task, error)
	GetAllTasksByUserAndStatus(
		ctx context.Context,
		userID uuid.UUID,
		status TaskStatus,
	) ([]*Task, error)
}

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user *User) (uuid.UUID, error)
	GetUsersByRole(ctx context.Context, role UserRoleType) ([]*User, error)
}
