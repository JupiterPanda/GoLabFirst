package libgrpc

import (
	"context"
	"goproject/internal/models"
	"goproject/protos/gen"
)

// GRPCServer реализует интерфейс LibraryServer из proto
type GRPCServer struct {
	gen.UnimplementedLibraryServer
	useCase UseCase
}

func NewGRPCServer(useCase UseCase) *GRPCServer {
	return &GRPCServer{useCase: useCase}
}

type UseCase interface {
	GetReaderBooksSepGoodAndBad(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error)
	RentBookByTitleAndReaderName(ctx context.Context, name, title string) error
	ReturnBookByTitleAndReaderName(ctx context.Context, name, title string) error

	GetAllBooks(ctx context.Context) ([]models.Book, error)
	GetBookByTitle(ctx context.Context, title string) (models.Book, error)
	GetBookIdByTitle(ctx context.Context, title string) (int, error)
	GetBookByID(ctx context.Context, id int) (models.Book, error)
	CreateBook(ctx context.Context, book models.Book) error
	DeleteBook(ctx context.Context, id int) error
	CheckCopiesOfBookByID(ctx context.Context, id int) error
	SubtractCopyOfBookById(ctx context.Context, id int) error
	AddCopyOfBookById(ctx context.Context, id int) error

	CreateBookInUse(ctx context.Context, readerId int, bookId int) error
	GetAllBooksInUse(ctx context.Context) ([]models.BookInUse, error)
	CountBookInUseByReaderId(ctx context.Context, readerId int) (int, error)
	GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error)
	GetBooksInUseByReaderId(ctx context.Context, readerId int) ([]models.BookInUse, error)
	DeleteBookInUse(ctx context.Context, readerId int, bookId int) error

	GetAllReaders(ctx context.Context) ([]models.Reader, error)
	CreateReader(ctx context.Context, reader models.Reader) error
	GetReaderIdByName(ctx context.Context, name string) (int, error)
	DeleteReader(ctx context.Context, id int) error
	UpdateReaderContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error
}
