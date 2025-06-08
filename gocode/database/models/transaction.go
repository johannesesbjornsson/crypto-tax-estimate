package models

import (
	"time"
)

type BaseTransaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Date        time.Time `gorm:"not null" json:"date"`
	Type        string    `gorm:"type:enum('income','lost','buy','sell');not null" json:"type"`
	Description string    `gorm:"type:varchar(100)" json:"description,omitempty"`
	Source      string    `gorm:"not null" json:"source"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Asset       string    `gorm:"not null" json:"asset"`
	ExternalID  string    `gorm:"type:varchar(100)" json:"external_id,omitempty"`
	UserID      uint      `gorm:"not null" json:"user_id"`
}

type TradeTransaction struct {
	Type          string  `gorm:"type:enum('buy','sell');not null" json:"type"`
	Price         float64 `gorm:"not null" json:"price"`
	QuoteCurrency string  `gorm:"not null" json:"quote_currency"`
	BaseTransaction
}

type SimpleTransaction struct {
	Type string `gorm:"type:enum('income','lost');not null" json:"type"`
	BaseTransaction
}
