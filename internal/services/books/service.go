package books

import (
	"context"
	"goproject/internal/models"
)

type Service struct {
	repo booksRepo
}

type booksRepo interface {
	GetByTitle(ctx context.Context, title string) (*models.Book, error)
	GetAll(ctx context.Context) ([]models.Book, error)
	CheckCopiesByID(ctx context.Context, id int) error
	CheckCopies(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, book *models.Book) error
	PlusCopyById(ctx context.Context, book *models.Book) error
	MinusCopyById(ctx context.Context, book *models.Book) error
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id int) (*models.Book, error)
	GetIdByTitle(ctx context.Context, title string) (int, error)
}

// NewService Конструктор сервиса книги
func NewService(repo booksRepo) *Service {
	return &Service{repo: repo}
}
