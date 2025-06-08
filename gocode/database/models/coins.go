package models

import (
	"time"
)

type Currency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	Type      string    `gorm:"type:enum('crypto','fiat','stablecoin');not null" json:"type"`
	PeggedTo  *uint     `gorm:"column:pegged_to" json:"pegged_to,omitempty"` // Nullable FK to currencies(id)
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
}