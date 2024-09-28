package migrations

import (
	"em-test/utils" // Импортируем кастомный логгер

	"github.com/golang-migrate/migrate/v4"                     // Импортируем пакет для миграций
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Импортируем драйвер для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Импортируем источник миграций из файлов
)

// MigrateUp выполняет миграции вверх до последней версии
func MigrateUp(dsn string) {
	// Инициализируем миграцию, указывая путь к миграциям и строку подключения (DSN)
	m, err := migrate.New(
		"file://db/migrations", // Путь к файлам миграций
		dsn)                    // Строка подключения к базе данных
	if err != nil {
		// Если произошла ошибка при инициализации миграции, логируем её и завершаем выполнение
		utils.ErrorLogger.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Применяем миграции вверх до последней версии
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// Если произошла ошибка при применении миграций, логируем её и завершаем выполнение
		utils.ErrorLogger.Fatalf("Не удалось применить миграции: %v", err)
	}

	// Логируем успешное применение миграций
	utils.InfoLogger.Println("Миграции успешно применены")
}

// MigrateDown откатывает миграции до начального состояния
func MigrateDown(dsn string) {
	// Инициализируем миграцию, указывая путь к миграциям и строку подключения (DSN)
	m, err := migrate.New(
		"file://db/migrations", // Путь к файлам миграций
		dsn)                    // Строка подключения к базе данных
	if err != nil {
		// Если произошла ошибка при инициализации миграции, логируем её и завершаем выполнение
		utils.ErrorLogger.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Откатываем миграции до начального состояния
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		// Если произошла ошибка при откате миграций, логируем её и завершаем выполнение
		utils.ErrorLogger.Fatalf("Не удалось откатить миграции: %v", err)
	}

	// Логируем успешный откат миграций
	utils.InfoLogger.Println("Миграции успешно откатаны")
}
