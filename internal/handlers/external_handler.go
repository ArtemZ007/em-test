package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ArtemZ007/em-test/internal/models"
)

// FetchSongLyrics calls the external API to enrich song details.
func FetchSongLyrics(group, title string) (*models.Song, error) {
	url := fmt.Sprintf("http://external-api-url/info?group=%s&song=%s", group, title)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch song details")
	}

	var enrichedSong models.Song
	if err := json.NewDecoder(resp.Body).Decode(&enrichedSong); err != nil {
		return nil, err
	}

	return &enrichedSong, nil
}
