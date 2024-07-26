package models

import "github.com/google/uuid"

type Admin struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Role   string    `gorm:"type:varchar(255)" validate:"oneof='SUPERUSER' 'ADMIN' 'EDITOR'" json:"role"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
}
