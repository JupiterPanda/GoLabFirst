package bookinuse

import (
	repository "goproject/internal/repositories/booksinuse"
)

// Service реализует доступ к данным по книгам
type Service struct {
	repo *repository.Repository
}

// NewService создаёт новый репозиторий с пулом соединений к базе
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
