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
	Gender       string    `json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	UserName     string    `gorm:"unique" json:"user_name"`
	UserPassword string    `gorm:"type:varchar(255)" json:"-"`
	MobileNo     string    `json:"mobile_no"`
	EmailID      string    `gorm:"unique" json:"email_id"`
	Files        []File    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // One-to-many relationship with Files
}

// File represents the file entity with additional information
type File struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Type      string    `gorm:"type:varchar(100)" validate:"oneof='profile' 'cover' 'document' "`
	FileName  string    `gorm:"type:text"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`   // Foreign key to the User table
	Size      int64     `gorm:"not null"`          // File size in bytes
	MimeType  string    `gorm:"size:100"`          // MIME type of the file
	CreatedAt time.Time `gorm:"autoCreateTime"`    // Creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime"`    // Update timestamp
	Path      string    `gorm:"size:500;not null"` // File path or location
	IsPublic  bool      `gorm:"default:false"`     // Public or private flag
}

type UserUiModel struct {
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birth_date"`
	UserName         string    `json:"user_name"`
	UserPassword     string    `json:"user_password"`
	MobileNo         string    `json:"mobile_no"`
	EmailID          string    `json:"email_id"`
	ProfilePicStatus bool      `json:"profile_pic_status"` // this will update whether profile pic will be publicly available or not
}
