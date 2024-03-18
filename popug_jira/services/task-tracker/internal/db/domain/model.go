package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskOpened       TaskStatus = "opened"
	TaskCompleted    TaskStatus = "completed"
	TaskUnknowStatus            = "unknown"
)

type Task struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	PublicID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null;unique"`

	UserID uuid.UUID `gorm:"type:uuid"`

	JiraId      string `gorm:"not null"`
	Description string `gorm:"not null"`

	Status TaskStatus `gorm:"not null;default:opened"`

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
	ID uuid.UUID `gorm:"type:uuid;not null;unique"`

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Username string       `gorm:"not null;unique"`
	Email    string       `gorm:"not null;unique"`
	Role     UserRoleType `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (User) TableName() string {
	return "users"
}
