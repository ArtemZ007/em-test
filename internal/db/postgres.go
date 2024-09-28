package db

import (
	"em-test/config" // Импортируем пакет с конфигурацией
	"fmt"            // Импортируем пакет для форматирования строк
	"log"            // Импортируем пакет для логирования
	"strconv"        // Импортируем пакет для преобразования типов
	"time"           // Импортируем пакет для работы со временем

	"gorm.io/driver/postgres" // Импортируем драйвер для PostgreSQL
	"gorm.io/gorm"            // Импортируем GORM - ORM библиотеку для Go
)

// Database структура для работы с БД
type Database struct {
	DB *gorm.DB // Поле для хранения подключения к базе данных
}

// InitDB инициализирует подключение к базе данных с использованием конфигурации
func InitDB(cfg *config.Config) (*Database, error) {
	portStr := strconv.Itoa(cfg.DBPort)

	// Формируем строку подключения (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, portStr, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *gorm.DB
	var err error

	// Пытаемся подключиться 5 раз с интервалом в 5 секунд
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Если подключение успешно, выходим из цикла
			break
		}
		// Если подключение не удалось, ждем 5 секунд и повторяем попытку
		log.Printf("Не удалось подключиться к базе данных. Попытка %d из 5. Ошибка: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		// Если все попытки не удались, возвращаем ошибку
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	// Получаем объект sql.DB из GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения SQL DB: %v", err)
	}

	// Настройка соединения
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	// Проверяем связь с базой данных
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при проверке связи с базой данных: %v", err)
	}

	log.Println("Подключение к базе данных успешно установлено")
	return &Database{DB: db}, nil
}
