package books

import (
	repository "goproject/internal/repositories/books"
)

type Service struct {
	repo *repository.Repository
}

// NewService Конструктор сервиса книги
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}
