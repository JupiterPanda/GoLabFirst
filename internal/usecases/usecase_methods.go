package usecases

import (
	"context"
	"fmt"
	"goproject/internal/models"
	constants "goproject/internal/package"
	"time"
)

// GetReaderBooksSepGoodAndBad Получить все книги у пользователя
func (u *UseCase) GetReaderBooksSepGoodAndBad(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error) {
	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		return nil, nil, fmt.Errorf("[useCase][[useCase][GetIdByTitle]] cannot get reader id by name: %w", err)
	}

	booksIds, err := u.bookInUseService.GetBooksInUseByReaderId(ctx, readerId)
	if err != nil {
		return nil, nil, fmt.Errorf("[useCase][GetBooksInUseByReaderId] cannot get books in use by reader id: %w", err)
	}

	var booksInUse []models.BookInUse
	for bookId, rentTime := range booksIds {
		bookInfo, err := u.bookService.GetByID(ctx, bookId)
		if err != nil {
			return nil, nil, fmt.Errorf("[useCase][GetByID] cannot get books in use by reader id: %w", err)
		}
		bookInUse := models.BookInUse{BookInfo: bookInfo, DateOfRent: rentTime}
		booksInUse = append(booksInUse, bookInUse)
	}

	var okBooks, badBooks []models.BookInUse
	for _, book := range booksInUse {
		if time.Since(book.DateOfRent) <= constants.TimeToExpire {
			okBooks = append(okBooks, book)
		} else {
			badBooks = append(badBooks, book)
		}
	}
	return okBooks, badBooks, nil
}

// RentBookByTitleAndReaderName Получить все книги у пользователя
func (u *UseCase) RentBookByTitleAndReaderName(ctx context.Context, name, title string) error {
	bookId, err := u.bookService.GetIdByTitle(ctx, title)
	if err != nil {
		return fmt.Errorf("[useCase][GetIdByTitle] cannot get book id by title: %w", err)
	}

	if err := u.bookService.CheckCopiesByID(ctx, bookId); err != nil {
		return fmt.Errorf("[useCase][CheckCopiesByID] no available copies: %w", err)
	}

	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		return fmt.Errorf("[useCase][GetIdByName] cannot get reader id by name: %w", err)
	}

	numOfBooks, err := u.bookInUseService.CountByReaderId(ctx, readerId)
	if err != nil {
		return fmt.Errorf("[useCase][CountByReaderId] cannot count books in use by reader: %w", err)
	}
	if numOfBooks >= 3 {
		return fmt.Errorf("max books limit reached")
	}

	if err := u.bookService.MinusCopyById(ctx, bookId); err != nil {
		return fmt.Errorf("[useCase][MinusCopyById] cannot decrease book copies: %w", err)
	}

	bookInfo, err := u.bookService.GetByID(ctx, bookId)
	if err != nil {
		return fmt.Errorf("[useCase][GetByID] cannot get book by id: %w", err)
	}

	bookInUse := models.BookInUse{BookInfo: bookInfo, DateOfRent: time.Now()}
	if err := u.bookInUseService.Create(ctx, bookInUse, readerId); err != nil {
		return fmt.Errorf("[useCase][Create] cannot create book in use: %w", err)
	}

	return nil
}

// ReturnBookByTitleAndReaderName Получить все книги у пользователя
func (u *UseCase) ReturnBookByTitleAndReaderName(ctx context.Context, name, title string) error {
	bookId, err := u.bookService.GetIdByTitle(ctx, title)
	if err != nil {
		return fmt.Errorf("[useCase][GetIdByTitle] cannot get book id by title: %w", err)
	}

	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		return fmt.Errorf("[useCase][GetIdByName] cannot get reader id by name: %w", err)
	}

	if err := u.bookService.PlusCopyById(ctx, bookId); err != nil {
		return fmt.Errorf("[useCase][PlusCopyById] cannot increase book copies: %w", err)
	}

	if err := u.bookInUseService.Delete(ctx, readerId, bookId); err != nil {
		return fmt.Errorf("[useCase][Delete] cannot delete book in use: %w", err)
	}

	return nil
}
