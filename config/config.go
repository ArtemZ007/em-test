package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// LoadConfig loads environment variables from .env file and returns a configured DB connection.
func LoadConfig() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}
