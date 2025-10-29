package libgrpc

import (
	"context"
	"goproject/internal/models"
	"goproject/protos/gen/librarypb"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func (s *GRPCServer) GetBookIdByTitle(ctx context.Context, req *librarypb.BookTitleRequest) (*librarypb.IdResponse, error) {
	bookID, err := s.useCase.GetBookIdByTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}
	res := &librarypb.IdResponse{Id: int32(bookID)}
	return res, nil
}

func (s *GRPCServer) GetBookByID(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.Book, error) {
	book, err := s.useCase.GetBookByID(ctx, int(req.Id))
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

func (s *GRPCServer) CreateBook(ctx context.Context, req *librarypb.CreateBookRequest) (*librarypb.MessageResponse, error) {
	book := models.Book{
		ID:     int(req.Book.ID),
		Title:  req.Book.Title,
		Copies: int(req.Book.Copies),
		Author: req.Book.Author,
		Issue:  req.Book.Issue.AsTime(),
	}
	err := s.useCase.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Book created")}
	return res, nil
}

func (s *GRPCServer) DeleteBook(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.DeleteBook(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Book deleted")}
	return res, nil
}

func (s *GRPCServer) CheckCopiesOfBookByID(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.CheckCopiesOfBookByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Alailable to rent")}
	return res, nil
}

func (s *GRPCServer) SubtractCopyOfBookById(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.SubtractCopyOfBookById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Copies decreased")}
	return res, nil
}

func (s *GRPCServer) AddCopyOfBookById(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.AddCopyOfBookById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Copies increased")}
	return res, nil
}
