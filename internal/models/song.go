package models

// Song представляет модель песни
type Song struct {
	ID          uint   `gorm:"primaryKey" json:"id"` // Уникальный идентификатор песни
	Group       string `json:"group"`                // Группа, исполняющая песню
	Song        string `json:"song"`                 // Название песни
	ReleaseDate string `json:"releaseDate"`          // Дата выпуска песни
	Text        string `json:"text"`                 // Текст песни
	Link        string `json:"link"`                 // Ссылка на песню
}
