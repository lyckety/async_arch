package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type HandleFuncPayoutToUser func(ctx context.Context, userPublicID uuid.UUID, paymentAmount int) error

type Repository interface {
	AccountingRepository
	TaskRepository
	UserRepository
	AccountRepository
	TransactionRepository
	BillingCyclesRepository
}

type AccountingRepository interface{}

type TaskRepository interface {
	CreateOrUpdateTask(context.Context, *Task) (uuid.UUID, error)
	GetTaskByPublicID(context.Context, uuid.UUID) (*Task, error)
}

type AccountRepository interface {
	GetUserBalance(context.Context, uuid.UUID) (int, error)
	GetTransactions(
		ctx context.Context,
		userPublicID uuid.UUID,
		startDate, endDate time.Time,
		transactionType TransactionType,
	) ([]TransactionWithDescriptionTask, error)
}

type UserRepository interface {
	CreateOrUpdateUser(context.Context, *User) (*User, error)
	GetUserByPublicID(context.Context, uuid.UUID) (*User, error)
}

type BillingCyclesRepository interface {
	CreateOrUpdateBillingCycleStatus(context.Context, BillingCycle) error
	CreateOpenedBillingCycleForUsers(ctx context.Context) error
	GetBillingCyclesByStatus(
		context.Context,
		BillingCycleStatus,
	) ([]BillingCycle, error)

	// Функция для процессинга выплаты заработка пользователю в рамках персонального биллинг цикла
	ProcessBillingCycle(
		context.Context,
		BillingCycle,
		HandleFuncPayoutToUser,
	) (*Transaction, error)
	AllStartedBillingCyclesToOpened(
		ctx context.Context,
	) error
}

type TransactionRepository interface {
	// for completed task
	CreateDebitTransaction(
		ctx context.Context,
		taskPublicID uuid.UUID,
		userPublicID uuid.UUID,
		time time.Time,
	) (*Transaction, error)
	// for assigned task
	CreateCreditTransaction(
		ctx context.Context,
		taskPublicID uuid.UUID,
		userPublicID uuid.UUID,
		time time.Time,
	) (*Transaction, error)
}
