package usecases

import (
	"context"
	"goproject/internal/models"
)

// GetAllReaders Получить все книги
func (u *UseCase) GetAllReaders(ctx context.Context) ([]models.Reader, error) {
	return u.readerService.GetAll(ctx)
}

// CreateReader Добавить читателя
func (u *UseCase) CreateReader(ctx context.Context, reader *models.Reader) error {
	return u.readerService.Create(ctx, reader)
}

// GetReaderIdByName Метод для получения ID читателя по имени
func (u *UseCase) GetReaderIdByName(ctx context.Context, name string) (int, error) {
	return u.readerService.GetIdByName(ctx, name)
}

// DeleteReader уладить читателя из readers
func (u *UseCase) DeleteReader(ctx context.Context, reader *models.Reader) error {
	return u.readerService.Delete(ctx, reader)
}

// UpdateReaderContactInfo обновит значения номера телефона или адреса
func (u *UseCase) UpdateReaderContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error {
	return u.readerService.UpdateContactInfo(ctx, readerId, phoneNumber, address)
}
