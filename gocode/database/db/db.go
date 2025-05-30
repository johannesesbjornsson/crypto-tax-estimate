package db

import (
	"database/sql"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDB() (*Database, *sql.DB) {
	log.Infof("Initializing Database connection")
	dsn := "user=personal host=localhost password=password dbname=cryptotax sslmode=disable"
	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL DB: %v", err)
	}

	return &Database{DB: dbInstance}, sqlDB
}

func (db *Database) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (db *Database) CreateOrUpdateUser(user *models.User) error {

	var existing models.User
	err := db.DB.Where("email = ?", user.Email).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing.ID > 0 {
		user.ID = existing.ID
		return db.DB.Save(user).Error
	}

	return db.DB.Create(user).Error
}
