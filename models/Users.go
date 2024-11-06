package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user model in the database
type User struct {
	ID           uuid.UUID     `gorm:"type:uuid;primary_key;" json:"id"`
	FindName     string        `gorm:"type:varchar(255)" json:"Find_name"`
	LastName     string        `gorm:"type:varchar(255)" json:"last_name"`
	Gender       string        `gorm:"type:varchar(255)" json:"gender"`
	BirthDate    time.Time     `json:"birth_date"`
	UserName     string        `gorm:"unique" json:"user_name"`
	UserPassword string        `gorm:"type:varchar(255)" json:"-"`
	MobileNo     string        `json:"mobile_no"`
	EmailID      string        `gorm:"unique" json:"email_id"`
	ProfileImage *ProfileImage `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID" json:"profile_image"` // One-to-many relationship with Files
	CoverImage   *CoverImage   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID" json:"cover_image"`   // One-to-many relationship with Files
}

// ProfileImage represents the profile image model in the database
type ProfileImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	FileName  string    `gorm:"type:text" json:"file_name"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`   // Foreign key to the User table
	Size      int64     `gorm:"not null" json:"size"`             // File size in bytes
	MimeType  string    `gorm:"size:100" json:"mime_type"`        // MIME type of the file
	Extension string    `gorm:"size:100" json:"extension"`        // File extension
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Update timestamp
	Path      string    `gorm:"size:500;not null" json:"path"`    // File path or location
	IsPublic  bool      `gorm:"default:true" json:"is_public"`    // Public or private flag
}

// CoverImage represents the cover image model in the database
type CoverImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	FileName  string    `gorm:"type:text" json:"file_name"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`   // Foreign key to the User table
	Size      int64     `gorm:"not null" json:"size"`             // File size in bytes
	MimeType  string    `gorm:"size:100" json:"mime_type"`        // MIME type of the file
	Extension string    `gorm:"size:100" json:"extension"`        // File extension
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Update timestamp
	Path      string    `gorm:"size:500;not null" json:"path"`    // File path or location
	IsPublic  bool      `gorm:"default:true" json:"is_public"`    // Public or private flag
}

type UserUiModel struct {
	FindName         string    `json:"Find_name"`
	LastName         string    `json:"last_name"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birth_date"`
	UserName         string    `json:"user_name"`
	UserPassword     string    `json:"user_password"`
	MobileNo         string    `json:"mobile_no"`
	EmailID          string    `json:"email_id"`
	ProfilePicStatus bool      `json:"profile_pic_status"` // this will update whether profile pic will be publicly available or not
	CoverPicStatus   bool      `json:"cover_pic_status"`   // this will update whether cover pic will be publicly available or notPicStatus bool      `json:"profile_pic_status"` // this will update whether profile pic will be publicly available or not
}
