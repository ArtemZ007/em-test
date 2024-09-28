package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"em-test/internal/db"
	"em-test/internal/models"

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

// Определяем константы для часто используемых значений
const (
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

// GetSongs godoc
// @Summary Получение списка песен
// @Description Возвращает список песен с поддержкой фильтрации и пагинации
// @Tags Песни
// @Accept  json
// @Produce  json
// @Param group query string false "Название группы"
// @Param limit query int false "Лимит записей на страницу"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {array} models.Song
// @Failure 500 {string} string "Ошибка на сервере"
// @Router /songs [get]
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
	query := h.db.DB.WithContext(ctx).Limit(limit).Offset(offset)
	if group != "" {
		query = query.Where("group = ?", group)
	}
	if err := query.Find(&songs).Error; err != nil {
		log.Printf("Не удалось получить песни: %v", err)
		http.Error(w, "Ошибка при получении песен", http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(songs)
}

// AddSong godoc
// @Summary Добавление новой песни
// @Description Создает новую песню и сохраняет в базе данных
// @Tags Песни
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Песня"
// @Success 201 {object} models.Song
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Ошибка на сервере"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	const invalidRequestFormat = "Неверный формат запроса"
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, invalidRequestFormat, http.StatusBadRequest)
		return
	}

	// Внешний запрос для получения данных о песне
	apiURL := fmt.Sprintf("http://external-api.com/info?group=%s&song=%s", song.Group, song.Song)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Ошибка при получении данных о песне", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var songDetail struct {
		ReleaseDate string `json:"release_date"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		http.Error(w, "Не удалось декодировать информацию о песне", http.StatusInternalServerError)
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	if err := h.db.DB.Create(&song).Error; err != nil {
		http.Error(w, "Не удалось сохранить песню", http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// GetSongHandler godoc
// @Summary Получение песни по ID
// @Description Возвращает информацию о песне по её ID
// @Tags Песни
// @Produce  json
// @Param id path string true "ID песни"
// @Success 200 {object} models.Song
// @Failure 404 {string} string "Песня не найдена"
// @Router /songs/{id} [get]
func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := h.db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(song)
}

// UpdateSongHandler godoc
// @Summary Обновление песни по ID
// @Description Обновляет информацию о песне по её ID
// @Tags Песни
// @Accept  json
// @Produce  json
// @Param id path string true "ID песни"
// @Param song body models.Song true "Данные для обновления"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Ошибка на сервере"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.db.DB.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		http.Error(w, "Ошибка при обновлении песни", http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(song)
}

// DeleteSongHandler godoc
// @Summary Удаление песни по ID
// @Description Удаляет песню из базы данных по её ID
// @Tags Песни
// @Param id path string true "ID песни"
// @Success 204 "Песня удалена"
// @Failure 500 {string} string "Ошибка на сервере"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := h.db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound)
		return
	}

	if err := h.db.DB.Delete(&song).Error; err != nil {
		http.Error(w, "Ошибка при удалении песни", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
