package utils

import (
	"log" // Импортируем пакет для логирования
	"os"  // Импортируем пакет для работы с файловой системой
)

// Инициализация логгера
var (
	InfoLogger  *log.Logger // Логгер для информационных сообщений
	ErrorLogger *log.Logger // Логгер для сообщений об ошибках
)

// InitLogger инициализирует логгеры для информационных и ошибочных сообщений
func InitLogger() {
	// Открываем файл для логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Если произошла ошибка при открытии файла, логируем её и завершаем выполнение
		log.Fatalf("Не удалось открыть файл для логирования: %v", err)
	}

	// Инициализируем логгеры
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)   // Логгер для информационных сообщений
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile) // Логгер для сообщений об ошибках
}

// Пример использования логгера для добавления новой песни
func LogNewSong(songName, groupName string) {
	// Логируем добавление новой песни с указанием её названия и группы
	InfoLogger.Printf("Добавление новой песни: %s от %s", songName, groupName)
}
