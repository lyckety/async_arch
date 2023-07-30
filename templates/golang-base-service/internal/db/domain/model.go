package domain

import "time"

type User struct {
	// UserID   uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`

	DateOfBirth      time.Time `gorm:"type:date"`
	RegistrationDate time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
