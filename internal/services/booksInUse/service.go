package bookinuse

import (
	"context"
	"goproject/internal/models"
	repository "goproject/internal/repositories/booksinuse"
)

// Service реализует доступ к данным по книгам
type Service struct {
	repo bookInUseRepo
}

type bookInUseRepo interface {
	Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error
	GetAll(ctx context.Context) ([]models.BookInUse, error)
}

// NewService создаёт новый репозиторий с пулом соединений к базе
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
