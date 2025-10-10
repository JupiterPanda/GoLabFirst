package books

import (
	"context"
	"errors"
	"goproject/internal/models"
)

// GetAll Получить все книги
func (s *Service) GetAll(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetAll(ctx)
}

// GetByTitle Получить книгу по названию
func (s *Service) GetByTitle(ctx context.Context, title string) (models.Book, error) {
	return s.repo.GetByTitle(ctx, title)
}

// GetIdByTitle Получить ID книги по названию
func (s *Service) GetIdByTitle(ctx context.Context, title string) (int, error) {
	return s.repo.GetIdByTitle(ctx, title)
}

// GetByID Получить книгу по ID
func (s *Service) GetByID(ctx context.Context, id int) (models.Book, error) {
	return s.repo.GetByID(ctx, id)
}

// Create Добавить новую книгу
func (s *Service) Create(ctx context.Context, book models.Book) error {
	if book.Title == "" || book.Author == "" || book.Copies < 1 {
		return errors.New("invalid book data")
	}
	return s.repo.Create(ctx, book)
}

// Delete удаляет из бд книгу (!!! Удалит книги и в таблице reader_books!!!)
func (s *Service) Delete(ctx context.Context, book models.Book) error {
	return s.repo.Delete(ctx, book)
}

// CheckCopiesByID проверяет кол-во книг в наличии по ID (if nil then copies > 0)
func (s *Service) CheckCopiesByID(ctx context.Context, id int) error {
	return s.repo.CheckCopiesByID(ctx, id)
}

// CheckCopies проверяет кол-во книг в наличии (if nil then copies > 0)
func (s *Service) CheckCopies(ctx context.Context, book models.Book) error {
	return s.repo.CheckCopies(ctx, book)
}

// MinusCopyById Уменьшить кол-во копий книги
func (s *Service) MinusCopyById(ctx context.Context, id int) error {
	return s.repo.MinusCopyById(ctx, id)
}

// PlusCopyById Увеличить кол-во копий книги
func (s *Service) PlusCopyById(ctx context.Context, id int) error {
	return s.repo.PlusCopyById(ctx, id)
}
