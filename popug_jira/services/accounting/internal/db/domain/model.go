package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	PublicID uuid.UUID `gorm:"type:uuid;not null;unique"`

	CostAssign   int `gorm:"not null"`
	CostComplete int `gorm:"not null"`

	Description string
	JiraId      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (Task) TableName() string {
	return "tasks"
}

type UserRoleType string

const (
	AdministratorRole UserRoleType = "administrator"
	BookkeepperRole   UserRoleType = "bookkeeper"
	ManagerRole       UserRoleType = "manager"
	WorkerRole        UserRoleType = "worker"
	UnknownRole       UserRoleType = "unknown"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	PublicID uuid.UUID `gorm:"type:uuid;not null;unique"`

	FirstName string
	LastName  string

	Username string
	Email    string
	Role     UserRoleType `gorm:"default:'unspecified'"`

	Balance int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (User) TableName() string {
	return "users"
}

type TransactionType string

const (
	TxTypeUnspecified TransactionType = "unspecified"
	TxTypeAssigned    TransactionType = "task_assigned"
	TxTypeCompleted   TransactionType = "task_completed"
	TxTypePaymented   TransactionType = "payment"
)

type Transaction struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	PublicID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null;unique"`
	UserPublicID uuid.UUID `gorm:"type:uuid;not null"`
	TaskPublicID uuid.UUID `gorm:"type:uuid;not null"`

	Debit    int             `gorm:"type:uuid;not null"`
	Credit   int             `gorm:"type:uuid;not null"`
	Type     TransactionType `gorm:"not null"`
	DateTime time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionWithDescriptionTask struct {
	Transaction
	TaskDescription string
	TaskJiraId      string
}

type BillingCycleStatus string

const (
	BillingCycleOpened    BillingCycleStatus = "opened"
	BillingCycleStarted   BillingCycleStatus = "started"
	BillingCycleCompleted BillingCycleStatus = "completed"
)

type BillingCycle struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserPublicID uuid.UUID `gorm:"type:uuid;not null"`

	StartDate time.Time
	EndDate   time.Time
	Status    BillingCycleStatus `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (BillingCycle) TableName() string {
	return "billing_cycles"
}
