package libgrpc

import (
	"context"
	"fmt"
	"goproject/internal/models"
	"goproject/protos/gen"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) GetAllBooks(ctx context.Context, _ *emptypb.Empty) (*gen.BooksResponse, error) {
	books, err := s.useCase.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}

	res := &gen.BooksResponse{}
	for _, book := range books {
		res.Books = append(res.Books, &gen.Book{
			ID:     int32(book.ID),
			Title:  book.Title,
			Copies: int32(book.Copies),
			Author: book.Author,
			Issue:  timestamppb.New(book.Issue),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetBookByTitle(ctx context.Context, req *gen.BookTitleRequest) (*gen.Book, error) {
	book, err := s.useCase.GetBookByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	res := &gen.Book{
		ID:     int32(book.ID),
		Title:  book.Title,
		Copies: int32(book.Copies),
		Author: book.Author,
		Issue:  timestamppb.New(book.Issue),
	}
	return res, nil
}

func (s *GRPCServer) GetBookIdByTitle(ctx context.Context, req *gen.BookTitleRequest) (*gen.IdResponse, error) {
	bookID, err := s.useCase.GetBookIdByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	res := &gen.IdResponse{Id: int32(bookID)}
	return res, nil
}

func (s *GRPCServer) GetBookByID(ctx context.Context, req *gen.BookIdRequest) (*gen.Book, error) {
	book, err := s.useCase.GetBookByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &gen.Book{
		ID:     int32(book.ID),
		Title:  book.Title,
		Copies: int32(book.Copies),
		Author: book.Author,
		Issue:  timestamppb.New(book.Issue),
	}
	return res, nil
}

func (s *GRPCServer) CreateBook(ctx context.Context, req *gen.CreateBookRequest) (*gen.MessageResponse, error) {
	if req == nil || req.Book == nil {
		return nil, fmt.Errorf("book is required")
	}
	fmt.Println(req.Book.Title)
	book := models.Book{
		Title:  req.Book.Title,
		Copies: int(req.Book.Copies),
		Author: req.Book.Author,
		Issue:  req.Book.Issue.AsTime(),
	}
	err := s.useCase.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book created"}
	return res, nil
}

func (s *GRPCServer) DeleteBook(ctx context.Context, req *gen.BookIdRequest) (*gen.MessageResponse, error) {
	err := s.useCase.DeleteBook(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book deleted"}
	return res, nil
}

func (s *GRPCServer) CheckCopiesOfBookByID(ctx context.Context, req *gen.BookIdRequest) (*gen.MessageResponse, error) {
	err := s.useCase.CheckCopiesOfBookByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Available to rent"}
	return res, nil
}

func (s *GRPCServer) SubtractCopyOfBookById(ctx context.Context, req *gen.BookIdRequest) (*gen.MessageResponse, error) {
	err := s.useCase.SubtractCopyOfBookById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Copies decreased"}
	return res, nil
}

func (s *GRPCServer) AddCopyOfBookById(ctx context.Context, req *gen.BookIdRequest) (*gen.MessageResponse, error) {
	err := s.useCase.AddCopyOfBookById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Copies increased"}
	return res, nil
}
