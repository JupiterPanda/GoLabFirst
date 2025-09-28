package handlers

import (
	"goproject/internal/services/books"
	"goproject/internal/services/readers"
)

type Handler struct {
	bookService   books.Service
	readerService readers.Service
}

func NewHandlers(
	bookService books.Service,
	readerService readers.Service,
) *Handler {
	return &Handler{
		bookService:   bookService,
		readerService: readerService,
	}
}
