package libgrpc

import (
	"context"
	"goproject/internal/handlers"
	"goproject/protos/gen/librarypb"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCServer реализует интерфейс LibraryServer из proto
type GRPCServer struct {
	librarypb.UnimplementedLibraryServer
	useCase handlers.UseCase
}

func NewGRPCServer(useCase handlers.UseCase) *GRPCServer {
	return &GRPCServer{useCase: useCase}
}

func (s *GRPCServer) GetAllBooks(ctx context.Context, _ *emptypb.Empty) (*librarypb.BooksResponse, error) {
	books, err := s.useCase.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}

	res := &librarypb.BooksResponse{}
	for _, book := range books {
		res.Books = append(res.Books, &librarypb.Book{
			ID:     int32(book.ID),
			Title:  book.Title,
			Copies: int32(book.Copies),
			Author: book.Author,
			Issue:  timestamppb.New(book.Issue),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetBookByTitle(ctx context.Context, req *librarypb.BookTitleRequest) (*librarypb.Book, error) {
	book, err := s.useCase.GetBookByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	res := &librarypb.Book{
		ID:     int32(book.ID),
		Title:  book.Title,
		Copies: int32(book.Copies),
		Author: book.Author,
		Issue:  timestamppb.New(book.Issue),
	}
	return res, nil
}
