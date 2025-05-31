package db

import (
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

func InitDB() *gorm.DB {
	dsn := "user=personal host=localhost password=password dbname=cryptotax sslmode=disable"
	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	return dbInstance
}