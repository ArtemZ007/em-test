package routes

import (
	"em-test/internal/handlers"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует все маршруты для API
func RegisterRoutes(router *mux.Router, songHandler *handlers.SongHandler) {
	const songIDRoute = "/songs/{id}"

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
