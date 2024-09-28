package db

import (
	"em-test/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database структура для работы с БД
type Database struct {
	DB *gorm.DB
}

// InitDB инициализирует подключение к базе данных с использованием конфигурации
func InitDB(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения SQL DB: %v", err)
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при проверке связи с базой данных: %v", err)
	}

	log.Println("Подключение к базе данных успешно установлено")
	return &Database{DB: db}, nil
}
