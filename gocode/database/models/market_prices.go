package models

import "time"

type MarketPrice struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	BaseCurrencyID  uint      `gorm:"not null" json:"base_currency_id"`
	QuoteCurrencyID uint      `gorm:"not null" json:"quote_currency_id"`
	Price           float64   `gorm:"not null" json:"price"`
	Timestamp       time.Time `gorm:"not null" json:"timestamp"`

	BaseCurrency  Currency `gorm:"foreignKey:BaseCurrencyID;references:ID" json:"base_currency"`
	QuoteCurrency Currency `gorm:"foreignKey:QuoteCurrencyID;references:ID" json:"quote_currency"`
}