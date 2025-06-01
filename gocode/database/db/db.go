package db

import (
	"database/sql"
	"fmt"
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

func (db *Database) GetTransactionsByEmail(email string, limit, offset int) ([]models.Transaction, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	query := db.DB.
		Where("user_id = ?", user.ID).
		Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *Database) CreateTransaction(tx *models.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	if tx.UserID == 0 {
		return fmt.Errorf("missing UserID on transaction")
	}

	return db.DB.Create(tx).Error
}

func (db *Database) GetFileUploadsByUserID(userID uint) ([]models.FileUploads, error) {
	if userID == 0 {
		return nil, fmt.Errorf("userID must be non-zero")
	}

	var uploads []models.FileUploads
	err := db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&uploads).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file uploads: %w", err)
	}

	return uploads, nil
}

func (db *Database) CreateFileUpload(fu *models.FileUploads) error {
	if fu == nil {
		return fmt.Errorf("file upload cannot be nil")
	}

	if fu.UserID == 0 {
		return fmt.Errorf("missing UserID on file upload")
	}

	return db.DB.Create(fu).Error
}