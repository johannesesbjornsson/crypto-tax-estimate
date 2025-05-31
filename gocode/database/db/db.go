package db

import (
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"

)

func InitDB() *gorm.DB {
	dsn := "user=personal host=localhost password=password dbname=cryptotax sslmode=disable"
	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	return dbInstance
}


func GetUserByEmail(email string) (*models.User, error) {
	db := InitDB()
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateOrUpdateUser(user *models.User) error {
	db := InitDB()

	var existing models.User
	err := db.Where("email = ?", user.Email).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing.ID > 0 {
		user.ID = existing.ID
		return db.Save(user).Error
	}

	return db.Create(user).Error
}