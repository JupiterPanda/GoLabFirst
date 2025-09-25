package bookinuse

import (
	"context"
	"goproject/internal/models"
	"time"
)

// Create добавляет в бд запись о новой книге у читателя
func (s *Service) Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error {
	return s.repo.Create(ctx, bookInUse, readerId)
}

// GetAll возвращает все книги из базы
func (s *Service) GetAll(ctx context.Context) ([]models.BookInUse, error) {
	return s.repo.GetAll(ctx)
}

// CountByReaderId ищет кол-во книг по ID клиента
func (s *Service) CountByReaderId(ctx context.Context, readerId int) (int, error) {
	return s.repo.CountByReaderId(ctx, readerId)
}

// GetReadersIdsByBookId ищет клиентов, взявших книгу по ID книги
func (s *Service) GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error) {
	return s.repo.GetReadersIdsByBookId(ctx, bookId)
}

// GetBooksInUseByReaderId ищет клиентов, взявших книгу по ID книги
func (s *Service) GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error) {
	return s.repo.GetBooksInUseByReaderId(ctx, readerId)
}

// Delete удаляет из бд запись об аренде книги
func (s *Service) Delete(ctx context.Context, readerId int, bookId int) error {
	return s.repo.Delete(ctx, readerId, bookId)
}
