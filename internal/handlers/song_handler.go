package handlers

import (
	"encoding/json" // Пакет для работы с JSON
	"fmt"           // Пакет для форматирования строк
	"net/http"      // Пакет для работы с HTTP
	"strconv"       // Пакет для конвертации строк в числа

	"em-test/internal/db"     // Пакет для работы с базой данных
	"em-test/internal/models" // Пакет с моделями данных
	"em-test/utils"           // Пакет с утилитами, включая логгер

	"github.com/gorilla/mux" // Пакет для работы с маршрутизацией
)

// SongHandler представляет обработчик для работы с песнями
type SongHandler struct {
	db *db.Database // Поле для работы с базой данных
}

// NewSongHandler создает новый экземпляр SongHandler
func NewSongHandler(db *db.Database) *SongHandler {
	return &SongHandler{db: db} // Возвращает новый экземпляр SongHandler
}

// Определяем константы для часто используемых значений
const (
	contentType     = "Content-Type"     // Константа для заголовка Content-Type
	applicationJSON = "application/json" // Константа для значения application/json
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
	ctx := r.Context()                               // Получаем контекст запроса
	var songs []models.Song                          // Создаем слайс для хранения песен
	filters := r.URL.Query()                         // Получаем параметры запроса
	limit, err := strconv.Atoi(filters.Get("limit")) // Преобразуем параметр limit в число
	if err != nil || limit <= 0 {
		limit = 10 // Значение по умолчанию
	}
	offset, err := strconv.Atoi(filters.Get("offset")) // Преобразуем параметр offset в число
	if err != nil || offset < 0 {
		offset = 0 // Значение по умолчанию
	}

	group := filters.Get("group")                                 // Получаем параметр group
	query := h.db.DB.WithContext(ctx).Limit(limit).Offset(offset) // Формируем запрос к базе данных с учетом лимита и смещения
	if group != "" {
		query = query.Where("group = ?", group) // Добавляем условие фильтрации по группе
	}
	if err := query.Find(&songs).Error; err != nil {
		utils.ErrorLogger.Printf("Не удалось получить песни: %v", err)              // Логируем ошибку
		http.Error(w, "Ошибка при получении песен", http.StatusInternalServerError) // Возвращаем ошибку клиенту
		return
	}

	w.Header().Set(contentType, applicationJSON) // Устанавливаем заголовок Content-Type
	json.NewEncoder(w).Encode(songs)             // Кодируем и отправляем список песен в формате JSON
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
	const invalidRequestFormat = "Неверный формат запроса" // Константа для сообщения об ошибке
	var song models.Song                                   // Создаем переменную для хранения данных песни
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, invalidRequestFormat, http.StatusBadRequest) // Возвращаем ошибку, если не удалось декодировать запрос
		return
	}

	// Внешний запрос для получения данных о песне
	apiURL := fmt.Sprintf("http://external-api.com/info?group=%s&song=%s", song.Group, song.Song)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Ошибка при получении данных о песне", http.StatusInternalServerError) // Возвращаем ошибку, если не удалось получить данные
		return
	}
	defer resp.Body.Close() // Закрываем тело ответа после завершения функции

	var songDetail struct {
		ReleaseDate string `json:"release_date"` // Поле для даты выпуска песни
		Text        string `json:"text"`         // Поле для текста песни
		Link        string `json:"link"`         // Поле для ссылки на песню
	}
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		http.Error(w, "Не удалось декодировать информацию о песне", http.StatusInternalServerError) // Возвращаем ошибку, если не удалось декодировать ответ
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate // Устанавливаем дату выпуска песни
	song.Text = songDetail.Text               // Устанавливаем текст песни
	song.Link = songDetail.Link               // Устанавливаем ссылку на песню

	if err := h.db.DB.Create(&song).Error; err != nil {
		http.Error(w, "Не удалось сохранить песню", http.StatusInternalServerError) // Возвращаем ошибку, если не удалось сохранить песню
		return
	}

	w.Header().Set(contentType, applicationJSON) // Устанавливаем заголовок Content-Type
	w.WriteHeader(http.StatusCreated)            // Устанавливаем статус ответа 201 Created
	json.NewEncoder(w).Encode(song)              // Кодируем и отправляем данные о песне в формате JSON
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
	vars := mux.Vars(r) // Получаем переменные пути из запроса
	id := vars["id"]    // Извлекаем ID песни

	var song models.Song // Создаем переменную для хранения данных песни
	if err := h.db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound) // Возвращаем ошибку, если песня не найдена
		return
	}

	w.Header().Set(contentType, applicationJSON) // Устанавливаем заголовок Content-Type
	json.NewEncoder(w).Encode(song)              // Кодируем и отправляем данные о песне в формате JSON
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
	vars := mux.Vars(r) // Получаем переменные пути из запроса
	id := vars["id"]    // Извлекаем ID песни

	var song models.Song // Создаем переменную для хранения данных песни
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest) // Возвращаем ошибку, если не удалось декодировать запрос
		return
	}

	if err := h.db.DB.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		http.Error(w, "Ошибка при обновлении песни", http.StatusInternalServerError) // Возвращаем ошибку, если не удалось обновить песню
		return
	}

	w.Header().Set(contentType, applicationJSON) // Устанавливаем заголовок Content-Type
	json.NewEncoder(w).Encode(song)              // Кодируем и отправляем обновленные данные о песне в формате JSON
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
	vars := mux.Vars(r) // Получаем переменные пути из запроса
	id := vars["id"]    // Извлекаем ID песни

	var song models.Song // Создаем переменную для хранения данных песни
	if err := h.db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Песня не найдена", http.StatusNotFound) // Возвращаем ошибку, если песня не найдена
		return
	}

	if err := h.db.DB.Delete(&song).Error; err != nil {
		http.Error(w, "Ошибка при удалении песни", http.StatusInternalServerError) // Возвращаем ошибку, если не удалось удалить песню
		return
	}

	w.WriteHeader(http.StatusNoContent) // Устанавливаем статус ответа 204 No Content
}
