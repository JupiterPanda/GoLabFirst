package models

import "time"

// BookInUse представляет книгу, взятую в аренду с датой
type BookInUse struct {
	BookInfo   Book      `json:"name_of_book"` // Книга
	DateOfRent time.Time `json:"date_of_rent"` // Дата взятия
}
