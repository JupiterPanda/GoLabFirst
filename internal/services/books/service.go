package books

import (
	"context"
	"goproject/internal/models"
	repository "goproject/internal/repositories/books"
)

type Service struct {
	repo  *repository.BookRepository
	books booksInterface
}

type booksInterface interface {
	GetAllBooks(context.Context) ([]models.Book, error)
	GetBookByTitle(ctx context.Context, title string) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) error
}

// New Конструктор сервиса книги
func New(repo *repository.BookRepository) *Service {
	return &Service{repo: repo}
}
