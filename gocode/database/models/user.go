package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Name         string    `gorm:"not null" json:"name"`
	TaxStartDate time.Time `json:"tax_start_date"`
	TaxEndDate   time.Time `json:"tax_end_date"`
	Currency     string    `gorm:"type:char(3);default:'USD'" json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}