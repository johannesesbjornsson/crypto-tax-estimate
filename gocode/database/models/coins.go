package models

import (
	"time"
)

type Currencies struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `gorm:"not null" json:"name"`
	BaseCurrency  bool      `gorm:"not null;default:false" json:"base_currency"`
}
