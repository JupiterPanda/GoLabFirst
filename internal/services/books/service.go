package books

import (
	"context"
	"errors"
	"goproject/internal/models"
	repository "goproject/internal/repositories/books"
)

type Service struct {
	repo *repository.BookRepository
}

// NewService Конструктор сервиса книги
func NewService(repo *repository.BookRepository) *Service {
	return &Service{repo: repo}
}

// GetAllBooks Получить все книги
func (s *Service) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetAll(ctx)
}

// CreateBook Добавить новую книгу
func (s *Service) CreateBook(ctx context.Context, book *models.Book) error {
	if book.Title == "" || book.Author == "" || book.Copies < 1 {
		return errors.New("invalid book data")
	}
	return s.repo.Create(ctx, book)
}

// GetBookByTitle Получить книгу по названию
func (s *Service) GetBookByTitle(ctx context.Context, title string) (*models.Book, error) {
	return s.repo.GetByTitle(ctx, title)
}
