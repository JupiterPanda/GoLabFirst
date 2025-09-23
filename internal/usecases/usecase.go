package usecases

import (
	"context"
	"goproject/internal/models"
	"time"
)

type UseCase struct {
	bookService      bookService
	readerService    readerService
	bookInUseService bookInUseService
}

type bookService interface {
	GetByTitle(ctx context.Context, title string) (*models.Book, error)
	GetAll(ctx context.Context) ([]models.Book, error)
	CheckCopiesByID(ctx context.Context, id int) error
	CheckCopies(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, book *models.Book) error
	PlusCopyById(ctx context.Context, book *models.Book) error
	MinusCopyById(ctx context.Context, id int) error //тут был косяк
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id int) (*models.Book, error)
	GetIdByTitle(ctx context.Context, title string) (int, error)
}

type readerService interface {
	GetAll(ctx context.Context) ([]models.Reader, error)
	Create(ctx context.Context, reader *models.Reader) error
	GetIdByName(ctx context.Context, name string) (int, error)
	Delete(ctx context.Context, readerId int) error
	UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error
}

type bookInUseService interface {
	Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error
	GetAll(ctx context.Context) ([]models.BookInUse, error)
	CountByReaderId(ctx context.Context, readerId int) (int, error)
	GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error)
	Delete(ctx context.Context, readerId int, bookId int) error
	GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error)
}

func New(bookService bookService, readerService readerService, bookInUseService bookInUseService) *UseCase {
	return &UseCase{bookService, readerService, bookInUseService}
}
