package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ArtemZ007/me-test/internal/db"
	"github.com/ArtemZ007/me-test/internal/models"

	"github.com/gorilla/mux"
)

// SongHandler представляет обработчик для работы с песнями
type SongHandler struct {
	db *db.Database
}

// NewSongHandler создает новый экземпляр SongHandler
func NewSongHandler(db *db.Database) *SongHandler {
	return &SongHandler{db: db}
}

// GetSongs обрабатывает запрос на получение списка песен с фильтрацией и пагинацией
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var songs []models.Song
	filters := r.URL.Query()
	limit, err := strconv.Atoi(filters.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // значение по умолчанию
	}
	offset, err := strconv.Atoi(filters.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // значение по умолчанию
	}

	group := filters.Get("group")
	if err := h.db.WithContext(ctx).Where("group = ?", group).Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		log.Printf("Не удалось получить песни: %v", err)
		http.Error(w, "Не удалось получить песни", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

// AddSong обрабатывает запрос на добавление новой песни
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	const invalidRequestFormat = "Неверный формат запроса"
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, invalidRequestFormat, http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("http://external-api.com/info?group=%s&song=%s", song.Group, song.Song)
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Не удалось получить информацию о песне", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		http.Error(w, "Не удалось декодировать информацию о песне", http.StatusInternalServerError)
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	if err := h.db.Create(&song).Error; err != nil {
		http.Error(w, "Не удалось сохранить песню", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// CreateSongHandler обрабатывает запрос на создание новой песни
func (h *SongHandler) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.db.Create(&song).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// GetSongsHandler обрабатывает запрос на получение всех песен
func (h *SongHandler) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := h.db.GetAllSongs()
	if err != nil {
		http.Error(w, "Не удалось получить песни", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

// GetSongHandler обрабатывает запрос на получение песни по ID
func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	song, err := h.db.GetSongByID(id)
	if err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(song)
}

// UpdateSongHandler обрабатывает запрос на обновление песни по ID
func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateSong(id, &song).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(song)
}

// DeleteSongHandler обрабатывает запрос на удаление песни по ID
func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.db.DeleteSong(id).Error; err != nil {
		http.Error(w, "Не удалось удалить песню", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
