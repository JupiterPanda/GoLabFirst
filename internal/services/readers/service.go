package readers

import (
	"context"
	"goproject/internal/models"
)

type Service struct {
	repo readerRepo
}

type readerRepo interface {
	Create(ctx context.Context, reader *models.Reader) error
	GetAll(ctx context.Context) ([]models.Reader, error)
	GetIdByName(ctx context.Context, name string) (int, error)
	Delete(ctx context.Context, reader *models.Reader) error
	UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error
}

// NewService Конструктор сервиса читателя
func NewService(repo readerRepo) *Service {
	return &Service{repo: repo}
}
