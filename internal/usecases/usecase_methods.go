package usecases

import (
	"context"
	"goproject/internal/models"
	constants "goproject/internal/package"
	"log"
	"time"
)

// GetReaderBooksSepGoodAndBad Получить все книги у пользователя
func (u *UseCase) GetReaderBooksSepGoodAndBad(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error) {
	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	//booksIds - мап(айди арендованных книг клиента и дата их аренды)
	booksIds, err := u.bookInUseService.GetBooksInUseByReaderId(ctx, readerId)
	if err != nil {
		log.Printf("unknown error")
	}

	var booksInUse []models.BookInUse
	for bookId, rentTime := range booksIds {
		bookInfo, err := u.bookService.GetByID(ctx, bookId)
		if err != nil {
			log.Printf("unknown error")
			continue
		}
		bookInUse := models.BookInUse{BookInfo: *bookInfo, DateOfRent: rentTime}
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
		log.Printf("book not found")
	}

	err = u.bookService.CheckCopiesByID(ctx, bookId)
	if err != nil {
		log.Printf("book not found")
	}

	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	numOfBooksOfReader, err := u.bookInUseService.CountByReaderId(ctx, readerId)
	if numOfBooksOfReader >= 3 {
		log.Printf("too many books")
	}
	if err != nil {
		log.Printf("reader not found")
	}

	err = u.bookService.MinusCopyById(ctx, bookId)
	if err != nil {
		log.Printf("reader not found")
	}

	bookInfo, err := u.bookService.GetByID(ctx, bookId)
	if err != nil {
		log.Printf("reader not found")
		return err
	}

	bookInUse := models.BookInUse{BookInfo: *bookInfo, DateOfRent: time.Now()}
	err = u.bookInUseService.Create(ctx, &bookInUse, readerId)
	if err != nil {
		log.Printf("reader not found")
	}

	return nil
}

// ReturnBookByTitleAndReaderName Получить все книги у пользователя
func (u *UseCase) ReturnBookByTitleAndReaderName(ctx context.Context, name, title string) error {

	bookId, err := u.bookService.GetIdByTitle(ctx, title)
	if err != nil {
		log.Printf("book not found")
	}

	err = u.bookService.CheckCopiesByID(ctx, bookId)
	if err != nil {
		log.Printf("book not found")
	}

	readerId, err := u.readerService.GetIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	err = u.bookService.PlusCopyById(ctx, bookId)
	if err != nil {
		log.Printf("reader not found")
	}

	err = u.bookInUseService.Delete(ctx, readerId, bookId)
	if err != nil {
		log.Printf("reader not found")
	}

	return nil
}
