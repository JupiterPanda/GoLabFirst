package usecases

import (
	"context"
	"errors"
	"goproject/internal/models"
)

// GetAllBooks Получить все книги
func (u *UseCase) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return u.bookService.GetAll(ctx)
}

// GetBookByTitle Получить книгу по названию
func (u *UseCase) GetBookByTitle(ctx context.Context, title string) (models.Book, error) {
	return u.bookService.GetByTitle(ctx, title)
}

// GetBookIdByTitle Получить ID книги по названию
func (u *UseCase) GetBookIdByTitle(ctx context.Context, title string) (int, error) {
	return u.bookService.GetIdByTitle(ctx, title)
}

// GetBookByID Получить книгу по ID
func (u *UseCase) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	return u.bookService.GetByID(ctx, id)
}

// CreateBook Добавить новую книгу
func (u *UseCase) CreateBook(ctx context.Context, book models.Book) error {
	if book.Title == "" {
		return errors.New("title invalid book data")
	}
	if book.Copies < 1 {
		return errors.New("copies invalid book data")
	}
	if book.Author == "" {
		return errors.New("author invalid book data")
	}
	return u.bookService.Create(ctx, book)
}

// DeleteBook удаляет из бд книгу (!!! Удалит книги и в таблице reader_books!!!)
func (u *UseCase) DeleteBook(ctx context.Context, id int) error {
	return u.bookService.Delete(ctx, id)
}

// CheckCopiesOfBookByID проверяет кол-во книг в наличии по ID (if nil then copies > 0)
func (u *UseCase) CheckCopiesOfBookByID(ctx context.Context, id int) error {
	return u.bookService.CheckCopiesByID(ctx, id)
}

// SubtractCopyOfBookById Уменьшить кол-во копий книги
func (u *UseCase) SubtractCopyOfBookById(ctx context.Context, id int) error {
	return u.bookService.SubtractCopyById(ctx, id)
}

// AddCopyOfBookById Увеличить кол-во копий книги
func (u *UseCase) AddCopyOfBookById(ctx context.Context, id int) error {
	return u.bookService.AddCopyById(ctx, id)
}
