package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	Name           string    `gorm:"not null" json:"name"`
	TaxStartDay    int       `gorm:"check:taxStartDay >= 1 AND taxStartDay <= 31" json:"taxStartDay"`
	TaxStartMonth  int       `gorm:"check:taxStartMonth >= 1 AND taxStartMonth <= 12" json:"taxStartMonth"`
	TaxEndDay      int       `gorm:"check:taxEndDay >= 1 AND taxEndDay <= 31" json:"taxEndDay"`
	TaxEndMonth    int       `gorm:"check:taxEndMonth >= 1 AND taxEndMonth <= 12" json:"taxEndMonth"`
	Currency       string    `gorm:"type:char(3);default:'USD'" json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}