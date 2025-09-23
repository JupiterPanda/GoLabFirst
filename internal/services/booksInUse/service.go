package bookinuse

import (
	"context"
	"goproject/internal/models"
	repository "goproject/internal/repositories/booksinuse"
	"time"
)

// Service реализует доступ к данным по книгам
type Service struct {
	repo bookInUseRepo
}

type bookInUseRepo interface {
	Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error
	GetAll(ctx context.Context) ([]models.BookInUse, error)
	CountByReaderId(ctx context.Context, readerId int) (int, error)
	GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error)
	Delete(ctx context.Context, readerId int, bookId int) error
	GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error)
}

// NewService создаёт новый репозиторий с пулом соединений к базе
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
