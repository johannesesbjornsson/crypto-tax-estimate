package models

import (
	"time"
)

type FileUploads struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:varchar(100)" json:"description,omitempty"`
	UserID      uint      `gorm:"not null" json:"user_id"`
}
