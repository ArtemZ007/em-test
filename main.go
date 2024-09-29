package main

import (
	"context"
	"em-test/config"            // Пакет для работы с конфигурацией
	"em-test/internal/db"       // Пакет для работы с базой данных
	"em-test/internal/handlers" // Пакет для обработки запросов
	"em-test/routes"            // Пакет для регистрации маршрутов
	"em-test/utils"             // Пакет для вспомогательных функций, таких как логирование
	"log"                       // Пакет для логирования
	"net/http"                  // Пакет для работы с HTTP сервером
	"os"                        // Пакет для работы с операционной системой
	"os/signal"                 // Для обработки системных сигналов (graceful shutdown)
	"time"                      // Для таймаута при завершении сервера

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
	// Здесь мы указываем путь к локальному файлу swagger.json или swagger.yaml, который находится в папке docs
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // Укажи правильный URL для Swagger документации
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

	// Запуск сервера с поддержкой graceful shutdown
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// Канал для системных сигналов (graceful shutdown)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Ожидание сигнала для завершения
	<-stopChan
	log.Println("Получен сигнал завершения, останавливаем сервер...")

	// Контекст с таймаутом для корректного завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Сервер завершил работу корректно.")
}
