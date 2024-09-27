package utils

import (
	"log"
	"os"
)

// Инициализация логгера
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// InitLogger инициализирует логгеры для информационных и ошибочных сообщений
func InitLogger() {
	// Открываем файл для логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логирования: %v", err)
	}

	// Инициализируем логгеры
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Пример использования логгера для добавления новой песни
func LogNewSong(songName, groupName string) {
	InfoLogger.Printf("Добавление новой песни: %s от %s", songName, groupName)
}
