package database

import (
	"em-test/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Выполнение миграций
	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
