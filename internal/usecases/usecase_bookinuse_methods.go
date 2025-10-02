package usecases

import (
	"context"
	"fmt"
	"goproject/internal/models"
)

// CreateBookInUse добавляет в бд запись о новой книге у читателя
func (u *UseCase) CreateBookInUse(ctx context.Context, bookInUse *models.BookInUse, readerId int, bookId int) error {
	bookPtr, err := u.bookService.GetByID(ctx, bookId)
	if err != nil {
		return fmt.Errorf("[useCase][CreateBookInUse] cannot get book by id: %w", err)
	}
	bookInUse.BookInfo = *bookPtr

	return u.bookInUseService.Create(ctx, bookInUse, readerId)
}

// GetAllBooksInUse возвращает все книги из базы
func (u *UseCase) GetAllBooksInUse(ctx context.Context) ([]models.BookInUse, error) {
	booksInUse, err := u.bookInUseService.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("[useCase][GetAllBooksInUse][GetAll] cannot get book in use: %w", err)
	}
	for i := range booksInUse {
		bookInfo, err := u.bookService.GetByID(ctx, booksInUse[i].BookInfo.ID)
		if err != nil {
			return nil, fmt.Errorf("[useCase][GetBooksInUseByReaderId][GetByID] cannot get book by id: %w", err)
		}
		booksInUse[i].BookInfo = *bookInfo
	}

	return booksInUse, err
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
func (u *UseCase) GetBooksInUseByReaderId(ctx context.Context, readerId int) ([]models.BookInUse, error) {
	mapOfBookIdAndDate, err := u.bookInUseService.GetBooksInUseByReaderId(ctx, readerId)
	if err != nil {
		return nil, fmt.Errorf("[useCase][GetBooksInUseByReaderId] cannot get book in use map: %w", err)
	}

	var booksInUse []models.BookInUse
	for bookId, date := range mapOfBookIdAndDate {
		bookInfo, err := u.bookService.GetByID(ctx, bookId)
		if err != nil {
			return nil, fmt.Errorf("[useCase][GetBooksInUseByReaderId][GetByID] cannot get book by id: %w", err)
		}
		//bookInUse := models.BookInUse{BookInfo: *bookInfo, DateOfRent: date}
		booksInUse = append(booksInUse, models.BookInUse{BookInfo: *bookInfo, DateOfRent: date})
	}
	return booksInUse, err
}

// DeleteBookInUse удаляет из бд запись об аренде книги
func (u *UseCase) DeleteBookInUse(ctx context.Context, readerId int, bookId int) error {
	return u.bookInUseService.Delete(ctx, readerId, bookId)
}
