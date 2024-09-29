package migrations

import (
	"em-test/config" // Подключаем конфигурацию
	"em-test/utils"  // Логирование
	"fmt"            // Форматирование строк
	"strconv"        // Для преобразования типов

	"github.com/golang-migrate/migrate/v4"                     // Пакет миграций
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Драйвер PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Источник миграций из файлов
)

// createDSN формирует строку подключения на основе конфигурации
func createDSN(cfg *config.Config) string {
	portStr := strconv.Itoa(cfg.DBPort)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, portStr, cfg.DBName)
}

// MigrateUp выполняет миграции вверх до последней версии
func MigrateUp(cfg *config.Config) {
	dsn := createDSN(cfg) // Используем функцию для получения строки подключения

	// Инициализируем миграцию
	m, err := migrate.New(
		"file://db/migrations", // Путь к файлам миграций
		dsn)
	if err != nil {
		utils.ErrorLogger.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Применяем миграции вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		utils.ErrorLogger.Fatalf("Не удалось применить миграции: %v", err)
	} else if err == migrate.ErrNoChange {
		utils.InfoLogger.Println("Нет новых миграций для применения")
	} else {
		utils.InfoLogger.Println("Миграции успешно применены")
	}
}

// MigrateDown откатывает миграции до начального состояния
func MigrateDown(cfg *config.Config) {
	dsn := createDSN(cfg) // Используем ту же функцию для DSN

	// Инициализируем миграцию
	m, err := migrate.New(
		"file://db/migrations",
		dsn)
	if err != nil {
		utils.ErrorLogger.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Откатываем миграции вниз
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		utils.ErrorLogger.Fatalf("Не удалось откатить миграции: %v", err)
	} else if err == migrate.ErrNoChange {
		utils.InfoLogger.Println("Нет миграций для отката")
	} else {
		utils.InfoLogger.Println("Миграции успешно откатаны")
	}
}
