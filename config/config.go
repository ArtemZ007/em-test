package config

import (
	"log"     // Пакет для логирования
	"os"      // Пакет для работы с операционной системой
	"strconv" // Пакет для конвертации строк

	"github.com/joho/godotenv" // Пакет для загрузки переменных окружения из файла .env
)

// Config представляет структуру конфигурации для подключения к базе данных
type Config struct {
	DBUser     string // Имя пользователя базы данных
	DBPassword string // Пароль пользователя базы данных
	DBHost     string // Хост базы данных
	DBPort     int    // Порт базы данных
	DBName     string // Имя базы данных
}

// LoadConfig загружает конфигурацию из файла .env и возвращает структуру Config
func LoadConfig() *Config {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env") // Логируем ошибку, если файл .env не найден или не может быть загружен
	}

	// Создаем структуру Config и заполняем её значениями из переменных окружения
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Некорректное значение порта в файле .env") // Логируем ошибку, если порт не является числом
	}

	config := Config{
		DBUser:     os.Getenv("DB_USER"),     // Получаем значение переменной окружения DB_USER
		DBPassword: os.Getenv("DB_PASSWORD"), // Получаем значение переменной окружения DB_PASSWORD
		DBHost:     os.Getenv("DB_HOST"),     // Получаем значение переменной окружения DB_HOST
		DBPort:     port,                     // Устанавливаем значение порта
		DBName:     os.Getenv("DB_NAME"),     // Получаем значение переменной окружения DB_NAME
	}

	// Проверяем, что все необходимые переменные окружения заданы
	if config.DBUser == "" || config.DBPassword == "" || config.DBHost == "" || config.DBName == "" {
		log.Fatal("Одно или несколько обязательных полей конфигурации не заданы в файле .env") // Логируем ошибку, если какое-либо поле не задано
	}

	// Возвращаем указатель на структуру Config
	return &config
}
