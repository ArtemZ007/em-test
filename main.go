package main

import (
	"em-test/api/routes"
	"em-test/internal/db"
	"log"
	"os"

	_ "github.com/ArtemZ007/em-test/docs" // путь к сгенерированной документации
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // библиотека для Swagger UI
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ArtemZ007/em-test/config"
	"github.com/ArtemZ007/em-test/utils"
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
	// Создаем новый роутер через Gin
	r := gin.Default()

	// Подключаем Swagger UI по адресу /swagger/*
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Загружаем конфигурацию из файла .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Инициализируем логгеры
	utils.InitLogger()

	// Инициализируем базу данных
	db.InitDB(cfg)

	// Регистрируем маршруты
	routes.RegisterRoutes(r)

	// Устанавливаем порт из переменной окружения или по умолчанию 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // значение по умолчанию
	}

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
