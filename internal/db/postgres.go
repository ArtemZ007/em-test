package db

import (
	"database/sql"
	"em-test/config"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Заменяем импорт на анонимный, если драйвер используется только для инициализации
)

// Database структура для работы с БД
type Database struct {
	DB *sql.DB
}

// InitDB инициализирует подключение к базе данных с использованием конфигурации
func InitDB(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при проверке связи с базой данных: %v", err)
	}

	log.Println("Подключение к базе данных успешно установлено")
	return &Database{DB: db}, nil
}
