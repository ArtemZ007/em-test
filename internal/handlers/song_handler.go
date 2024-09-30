package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtemZ007/em-test/internal/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SongHandler handles requests related to songs.
type SongHandler struct {
	DB *gorm.DB
}

// @Summary Get a list of songs
// @Description Retrieve a list of songs with optional filters (group, title) and pagination.
// @Tags Songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group"
// @Param title query string false "Filter by title"
// @Param limit query int false "Number of results to return"
// @Param offset query int false "Number of results to skip"
// @Success 200 {array} models.Song
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	group := r.URL.Query().Get("group")
	title := r.URL.Query().Get("title")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	query := h.DB.Limit(limit).Offset(offset)
	if group != "" {
		query = query.Where("group = ?", group)
	}
	if title != "" {
		query = query.Where("title = ?", title)
	}

	if result := query.Find(&songs); result.Error != nil {
		log.Println("Error fetching songs:", result.Error)
		http.Error(w, "Failed to fetch songs", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched %d songs", len(songs))
	json.NewEncoder(w).Encode(songs)
}

// @Summary Add a new song
// @Description Add a new song and enrich its data using an external API.
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song to add"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Adding song: %s - %s", song.Group, song.Title)

	// Enrich song data from external API
	enrichedSong, err := FetchSongLyrics(song.Group, song.Title)
	if err != nil {
		log.Println("Error fetching external API:", err)
		http.Error(w, "Failed to enrich song data", http.StatusInternalServerError)
		return
	}

	// Save enriched data
	song.Text = enrichedSong.Text
	song.ReleaseDate = time.Now() // Simulate external date fetching
	song.Link = enrichedSong.Link

	if result := h.DB.Create(&song); result.Error != nil {
		log.Println("Error saving song:", result.Error)
		http.Error(w, "Failed to add song", http.StatusInternalServerError)
		return
	}

	log.Printf("Song added successfully: %s - %s", song.Group, song.Title)
	json.NewEncoder(w).Encode(song)
}

// @Summary Delete a song
// @Description Delete a song by its ID.
// @Tags Songs
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if result := h.DB.Delete(&models.Song{}, id); result.Error != nil {
		log.Println("Error deleting song:", result.Error)
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	log.Printf("Song with ID %s deleted successfully", id)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update a song
// @Description Update a song's details.
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedSong models.Song

	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if result := h.DB.Model(&models.Song{}).Where("id = ?", id).Updates(&updatedSong); result.Error != nil {
		log.Println("Error updating song:", result.Error)
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	log.Printf("Song with ID %s updated successfully", id)
	json.NewEncoder(w).Encode(updatedSong)
}
