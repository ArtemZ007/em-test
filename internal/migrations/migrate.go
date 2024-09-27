package migrations

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp выполняет миграции вверх до последней версии
func MigrateUp(dsn string) {
	// Инициализируем миграцию
	m, err := migrate.New(
		"file://db/migrations",
		dsn)
	if err != nil {
		log.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Не удалось применить миграции: %v", err)
	}

	log.Println("Миграции успешно применены")
}

// MigrateDown откатывает миграции до начального состояния
func MigrateDown(dsn string) {
	// Инициализируем миграцию
	m, err := migrate.New(
		"file://db/migrations",
		dsn)
	if err != nil {
		log.Fatalf("Не удалось инициализировать миграцию: %v", err)
	}

	// Откатываем миграции
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Не удалось откатить миграции: %v", err)
	}

	log.Println("Миграции успешно откатаны")
}
