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

	// Чтение порта базы данных из переменной окружения
	portStr := os.Getenv("POSTGRES_PORT") // Считываем переменную окружения для порта как строку
	port, err := strconv.Atoi(portStr)    // Конвертируем строку в число
	if err != nil {
		log.Fatalf("Некорректное значение порта в файле .env: %v", err) // Логируем ошибку, если порт не является числом
	}

	// Создаем структуру Config и заполняем её значениями из переменных окружения
	config := Config{
		DBUser:     os.Getenv("POSTGRES_USER"),     // Получаем значение переменной окружения POSTGRES_USER
		DBPassword: os.Getenv("POSTGRES_PASSWORD"), // Получаем значение переменной окружения POSTGRES_PASSWORD
		DBHost:     os.Getenv("POSTGRES_HOST"),     // Получаем значение переменной окружения POSTGRES_HOST
		DBPort:     port,                           // Устанавливаем значение порта, уже преобразованное в int
		DBName:     os.Getenv("POSTGRES_DB"),       // Получаем значение переменной окружения POSTGRES_DB
	}

	// Проверяем, что все необходимые переменные окружения заданы
	if config.DBUser == "" || config.DBPassword == "" || config.DBHost == "" || config.DBName == "" {
		log.Fatal("Одно или несколько обязательных полей конфигурации не заданы в файле .env") // Логируем ошибку, если какое-либо поле не задано
	}

	// Логируем загруженные значения для отладки (можно убрать на продакшене)
	log.Printf("Конфигурация загружена: хост=%s порт=%d пользователь=%s база данных=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBName)

	// Возвращаем указатель на структуру Config
	return &config
}
