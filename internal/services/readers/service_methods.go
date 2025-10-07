package readers

import (
	"context"
	"goproject/internal/models"
)

// GetAll Получить все книги
func (s *Service) GetAll(ctx context.Context) ([]models.Reader, error) {
	return s.repo.GetAll(ctx)
}

// Create Добавить читателя
func (s *Service) Create(ctx context.Context, reader models.Reader) error {
	return s.repo.Create(ctx, reader)
}

// GetIdByName Метод для получения ID читателя по имени
func (s *Service) GetIdByName(ctx context.Context, name string) (int, error) {
	return s.repo.GetIdByName(ctx, name)
}

// Delete уладить читателя из readers
func (s *Service) Delete(ctx context.Context, reader models.Reader) error {
	return s.repo.Delete(ctx, reader)
}

// UpdateContactInfo обновит значения номера телефона или адреса
func (s *Service) UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error {
	return s.repo.UpdateContactInfo(ctx, readerId, phoneNumber, address)
}
