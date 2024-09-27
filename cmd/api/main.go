package main

import (
	"em-test/api/routes"
	"em-test/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/ArtemZ007/em-test/config"
	"github.com/ArtemZ007/em-test/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Инициализируем логгеры
	utils.InitLogger()

	// Инициализируем базу данных
	db.InitDB(cfg)

	// Создаем новый роутер
	r := mux.NewRouter()

	// Регистрируем маршруты
	routes.RegisterRoutes(r)

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // значение по умолчанию
	}
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
