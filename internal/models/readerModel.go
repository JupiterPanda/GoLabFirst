package models

import "time"

// Reader представляет пользователя библиотеки
type Reader struct {
	ID          int         `json:"ID"`            // ID читателя
	Name        string      `json:"name"`          // Имя читателя
	PhoneNumber string      `json:"phone_number"`  // Номер телефона
	Address     string      `json:"address"`       // Адрес
	DateOfBirth time.Time   `json:"date_of_birth"` // Дата рождения
	BooksInUse  []BookInUse `json:"books_in_use"`  // Книги, взятые в пользование
}

// BookInUse представляет книгу, взятую в аренду с датой
type BookInUse struct {
	NameOfBook Book      `json:"name_of_book"` // Книга
	DateOfRent time.Time `json:"date_of_rent"` // Дата взятия
}
