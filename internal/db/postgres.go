package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ArtemZ007/em-test/config"
	_ "github.com/lib/pq"
)

// DB представляет собой глобальную переменную для подключения к базе данных
var DB *sql.DB

// InitDB инициализирует подключение к базе данных с использованием конфигурации
func InitDB(config *config.Config) {
	var err error
	// Формируем строку подключения (DSN) на основе конфигурации
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	// Открываем подключение к базе данных
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Проверяем подключение к базе данных
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Не удалось выполнить ping к базе данных: %v", err)
	}

	fmt.Println("Успешное подключение к базе данных!")
}
