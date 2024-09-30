package models

import (
	"time"
)

// Song defines the structure for songs stored in the database.
type Song struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Group       string    `json:"group" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
