package main

import (
	"em-test/config"            // Пакет для работы с конфигурацией
	"em-test/internal/db"       // Пакет для работы с базой данных
	"em-test/internal/handlers" // Пакет для обработки запросов
	"em-test/routes"            // Пакет для регистрации маршрутов
	"em-test/utils"             // Пакет для вспомогательных функций, таких как логирование
	"log"                       // Пакет для логирования
	"net/http"                  // Пакет для работы с HTTP сервером
	"os"                        // Пакет для работы с операционной системой

	"github.com/gorilla/mux"                     // Пакет для маршрутизации HTTP запросов
	httpSwagger "github.com/swaggo/http-swagger" // Пакет для интеграции Swagger UI
)

// @title Music Library API
// @version 1.0
// @description This is a music library service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Создаем новый роутер через mux
	router := mux.NewRouter()

	// Подключаем Swagger UI по адресу /swagger/*
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажи правильный URL для Swagger документации
	))

	// Загружаем конфигурацию из файла .env
	cfg := config.LoadConfig()

	// Инициализируем логгеры
	utils.InitLogger()

	// Инициализируем базу данных
	dbInstance, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Создаем SongHandler
	songHandler := handlers.NewSongHandler(dbInstance)

	// Регистрируем маршруты
	routes.RegisterRoutes(router, songHandler)

	// Устанавливаем порт из переменной окружения или по умолчанию 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
