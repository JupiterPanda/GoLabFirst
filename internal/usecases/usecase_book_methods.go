package usecases

import (
	"context"
	"errors"
	"goproject/internal/models"
)

// GetAll Получить все книги
func (u *UseCase) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return u.bookService.GetAll(ctx)
}

// GetByTitle Получить книгу по названию
func (u *UseCase) GetBookByTitle(ctx context.Context, title string) (*models.Book, error) {
	return u.bookService.GetByTitle(ctx, title)
}

// GetIdByTitle Получить ID книги по названию
func (u *UseCase) GetBookIdByTitle(ctx context.Context, title string) (int, error) {
	return u.bookService.GetIdByTitle(ctx, title)
}

// GetByID Получить книгу по ID
func (u *UseCase) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	return u.bookService.GetByID(ctx, id)
}

// Create Добавить новую книгу
func (u *UseCase) CreateBook(ctx context.Context, book *models.Book) error {
	if book.Title == "" || book.Author == "" || book.Copies < 1 {
		return errors.New("invalid book data")
	}
	return u.bookService.Create(ctx, book)
}

// Delete удаляет из бд книгу (!!! Удалит книги и в таблице reader_books!!!)
func (u *UseCase) DeleteBook(ctx context.Context, book *models.Book) error {
	return u.bookService.Delete(ctx, book)
}

// CheckCopiesByID проверяет кол-во книг в наличии по ID (if nil then copies > 0)
func (u *UseCase) CheckCopiesOfBookByID(ctx context.Context, id int) error {
	return u.bookService.CheckCopiesByID(ctx, id)
}

// CheckCopies проверяет кол-во книг в наличии (if nil then copies > 0)
func (u *UseCase) CheckCopiesOfBook(ctx context.Context, book *models.Book) error {
	return u.bookService.CheckCopies(ctx, book)
}

// MinusCopyById Уменьшить кол-во копий книги
func (u *UseCase) MinusCopyOfBookById(ctx context.Context, id int) error {
	return u.bookService.MinusCopyById(ctx, id)
}

// PlusCopyById Увеличить кол-во копий книги
func (u *UseCase) PlusCopyOfBookById(ctx context.Context, id int) error {
	return u.bookService.PlusCopyById(ctx, id)
}
