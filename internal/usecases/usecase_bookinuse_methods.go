package usecases

import (
	"context"
	"goproject/internal/models"
	"time"
)

// CreateBookInUse добавляет в бд запись о новой книге у читателя
func (u *UseCase) CreateBookInUse(ctx context.Context, bookInUse *models.BookInUse, readerId int) error {
	return u.bookInUseService.Create(ctx, bookInUse, readerId)
}

// GetAllBooksInUse возвращает все книги из базы
func (u *UseCase) GetAllBooksInUse(ctx context.Context) ([]models.BookInUse, error) {
	return u.bookInUseService.GetAll(ctx)
}

// CountBookInUseByReaderId ищет кол-во книг по ID клиента
func (u *UseCase) CountBookInUseByReaderId(ctx context.Context, readerId int) (int, error) {
	return u.bookInUseService.CountByReaderId(ctx, readerId)
}

// GetReadersIdsByBookId возвращает срез ID читателей, арендовавших книгу по ID книги
func (u *UseCase) GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error) {
	return u.bookInUseService.GetReadersIdsByBookId(ctx, bookId)
}

// GetBooksInUseByReaderId ищет клиентов, взявших книгу по ID книги
func (u *UseCase) GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error) {
	return u.bookInUseService.GetBooksInUseByReaderId(ctx, readerId)
}

// DeleteBookInUse удаляет из бд запись об аренде книги
func (u *UseCase) DeleteBookInUse(ctx context.Context, readerId int, bookId int) error {
	return u.bookInUseService.Delete(ctx, readerId, bookId)
}
