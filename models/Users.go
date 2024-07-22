package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user model in the database
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	ProfilePhoto string    `json:"profile_photo"`
	Gender       string    `json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	UserName     string    `gorm:"unique" json:"user_name"`
	UserPassword string    `gorm:"type:varchar(255)" json:"-"`
	MobileNo     string    `json:"mobile_no"`
	EmailID      string    `gorm:"unique" json:"email_id"`
}
type UserUiModel struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Gender       string    `json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	UserName     string    `gorm:"unique" json:"user_name"`
	UserPassword string    `gorm:"type:varchar(255)" json:"user_password"`
	MobileNo     string    `json:"mobile_no"`
	EmailID      string    `gorm:"unique" json:"email_id"`
}
