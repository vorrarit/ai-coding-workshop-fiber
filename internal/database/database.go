package database

import (
	"fiber-api/internal/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(databasePath string) *Database {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Transfer{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return &Database{DB: db}
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}
