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

// Create Добавить новую книгу
func (s *Service) Create(ctx context.Context, book *models.Book) error {
	if book.Title == "" || book.Author == "" || book.Copies < 1 {
		return errors.New("invalid book data")
	}
	return s.repo.Create(ctx, book)
}

// GetByTitle Получить книгу по названию
func (s *Service) GetByTitle(ctx context.Context, title string) (*models.Book, error) {
	return s.repo.GetByTitle(ctx, title)
}

// Delete удаляет из бд книгу (!!! Удалит книги и в таблице reader_books!!!)
func (s *Service) Delete(ctx context.Context, book *models.Book) error {
	return s.repo.Delete(ctx, book)
}

func (s *Service) CheckCopiesByID(ctx context.Context, id int) error {
	return s.repo.CheckCopiesByID(ctx, id)
}

func (s *Service) CheckCopies(ctx context.Context, book *models.Book) error {
	return s.repo.CheckCopies(ctx, book)
}

func (s *Service) MinusCopyById(ctx context.Context, book *models.Book) error {
	return s.repo.MinusCopyById(ctx, book)
}

func (s *Service) PlusCopyById(ctx context.Context, book *models.Book) error {
	return s.repo.PlusCopyById(ctx, book)
}
