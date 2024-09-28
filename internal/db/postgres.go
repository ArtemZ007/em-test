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
	// Преобразуем порт в строку
	portStr := strconv.Itoa(cfg.DBPort)

	// Формируем строку подключения (DSN) с использованием параметров из конфигурации
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, portStr, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Открываем соединение с базой данных с использованием GORM и драйвера PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Если произошла ошибка при подключении, возвращаем её
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	// Получаем объект sql.DB из GORM для настройки соединения
	sqlDB, err := db.DB()
	if err != nil {
		// Если произошла ошибка при получении sql.DB, возвращаем её
		return nil, fmt.Errorf("ошибка получения SQL DB: %v", err)
	}

	// Устанавливаем максимальное время жизни соединения
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// Устанавливаем максимальное количество открытых соединений
	sqlDB.SetMaxOpenConns(20)
	// Устанавливаем максимальное количество соединений в пуле ожидания
	sqlDB.SetMaxIdleConns(10)

	// Проверяем связь с базой данных
	if err := sqlDB.Ping(); err != nil {
		// Если произошла ошибка при проверке связи, возвращаем её
		return nil, fmt.Errorf("ошибка при проверке связи с базой данных: %v", err)
	}

	// Логируем успешное подключение к базе данных
	log.Println("Подключение к базе данных успешно установлено")
	// Возвращаем объект Database с установленным подключением
	return &Database{DB: db}, nil
}
