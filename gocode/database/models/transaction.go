package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Date        time.Time `gorm:"not null" json:"date"`
	Description string    `gorm:"type:varchar(100)" json:"description,omitempty"`
	Venue       string    `gorm:"not null" json:"venue"`                         
	Source      string    `gorm:"not null" json:"source"`                         
	Type        string    `gorm:"type:enum('Income','Buy','Sell','Lost');not null" json:"type"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Asset       string    `gorm:"not null" json:"asset"`
	UserID      uint      `gorm:"not null" json:"user_id"`

}