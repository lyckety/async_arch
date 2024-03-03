package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleType string

const (
	Administrator UserRoleType = "administrator"
	Bookkeepper   UserRoleType = "bookkeeper"
	Manager       UserRoleType = "manager"
	Worker        UserRoleType = "worker"
	Unknown       UserRoleType = "unknown"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Username string       `gorm:"not null;unique"`
	Password string       `gorm:"not null"`
	Email    string       `gorm:"not null;unique"`
	Role     UserRoleType `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

func (User) TableName() string {
	return "users"
}
