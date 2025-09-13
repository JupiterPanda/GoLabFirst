package models

import "time"

type Book struct {
	ID     int       `json:"ID"`     // ID Книги
	Title  string    `json:"title"`  // Название книги
	Copies int       `json:"copies"` // Количество копий
	Author string    `json:"author"` // Автор книги
	Issue  time.Time `json:"issue"`  // Дата выпуска
}
