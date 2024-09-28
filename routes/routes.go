package routes

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует все маршруты для API
func RegisterRoutes(router *mux.Router) {
	// Маршрут для получения всех песен
	router.HandleFunc("/songs", controllers.GetSongs).Methods("GET")

	// Маршрут для создания новой песни
	router.HandleFunc("/songs", controllers.CreateSong).Methods("POST")

	// Маршрут для удаления песни по ID
	router.HandleFunc("/songs/{id}", controllers.DeleteSong).Methods("DELETE")

	// Маршрут для обновления песни по ID
	router.HandleFunc("/songs/{id}", controllers.UpdateSong).Methods("PUT")
}
