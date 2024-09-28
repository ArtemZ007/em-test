package routes

import (
	"em-test/internal/handlers" // Импортируем пакет с обработчиками запросов

	"github.com/gorilla/mux" // Импортируем пакет для маршрутизации
)

// RegisterRoutes регистрирует все маршруты для API
func RegisterRoutes(router *mux.Router, songHandler *handlers.SongHandler) {
	const songIDRoute = "/songs/{id}" // Константа для маршрута с ID песни

	// Маршрут для получения всех песен
	router.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")

	// Маршрут для создания новой песни
	router.HandleFunc("/songs", songHandler.AddSong).Methods("POST")

	// Маршрут для удаления песни по ID
	router.HandleFunc(songIDRoute, songHandler.DeleteSongHandler).Methods("DELETE")

	// Маршрут для обновления песни по ID
	router.HandleFunc(songIDRoute, songHandler.UpdateSongHandler).Methods("PUT")

	// Маршрут для получения песни по ID
	router.HandleFunc(songIDRoute, songHandler.GetSongHandler).Methods("GET")
}
