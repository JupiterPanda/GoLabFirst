package libgrpc

import (
	"context"
	"goproject/protos/gen"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) RentBookByTitleAndReaderName(ctx context.Context, req *gen.RentOrReturnBookRequest) (*gen.MessageResponse, error) {
	err := s.useCase.RentBookByTitleAndReaderName(ctx, req.Name, req.Title)
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book rented successfully"}
	return res, nil
}

func (s *GRPCServer) ReturnBookByTitleAndReaderName(ctx context.Context, req *gen.RentOrReturnBookRequest) (*gen.MessageResponse, error) {
	err := s.useCase.ReturnBookByTitleAndReaderName(ctx, req.Name, req.Title)
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book returned successfully"}
	return res, nil
}

func (s *GRPCServer) GetReaderBooksSepGoodAndBad(ctx context.Context, req *gen.ReaderNameRequest) (*gen.ReadersBooksResponse, error) {
	okBooks, badBooks, err := s.useCase.GetReaderBooksSepGoodAndBad(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	res := &gen.ReadersBooksResponse{}

	for _, bookInUse := range okBooks {
		res.OkBooks = append(res.OkBooks, &gen.BookInUse{
			BookInfo: &gen.Book{
				ID:     int32(bookInUse.BookInfo.ID),
				Title:  bookInUse.BookInfo.Title,
				Copies: int32(bookInUse.BookInfo.Copies),
				Author: bookInUse.BookInfo.Author,
				Issue:  timestamppb.New(bookInUse.BookInfo.Issue),
			},
			DateOfRent: timestamppb.New(bookInUse.DateOfRent),
		})
	}
	for _, bookInUse := range badBooks {
		res.BadBooks = append(res.BadBooks, &gen.BookInUse{
			BookInfo: &gen.Book{
				ID:     int32(bookInUse.BookInfo.ID),
				Title:  bookInUse.BookInfo.Title,
				Copies: int32(bookInUse.BookInfo.Copies),
				Author: bookInUse.BookInfo.Author,
				Issue:  timestamppb.New(bookInUse.BookInfo.Issue),
			},
			DateOfRent: timestamppb.New(bookInUse.DateOfRent),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetReadersIdsByBookId(ctx context.Context, req *gen.BookIdRequest) (*gen.ReaderIdsResponse, error) {
	listOfBooksInUseOfReader, err := s.useCase.GetReadersIdsByBookId(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	int32ListOfBooksInUseOfReader := make([]int32, len(listOfBooksInUseOfReader), 3)
	for i, id := range listOfBooksInUseOfReader {
		int32ListOfBooksInUseOfReader[i] = int32(id)
	}
	res := &gen.ReaderIdsResponse{ReaderIds: int32ListOfBooksInUseOfReader}
	return res, nil
}

func (s *GRPCServer) GetAllBooksInUse(ctx context.Context, _ *emptypb.Empty) (*gen.BooksInUseResponse, error) {
	booksInUse, err := s.useCase.GetAllBooksInUse(ctx)
	if err != nil {
		return nil, err
	}

	res := &gen.BooksInUseResponse{}
	for _, bookInUse := range booksInUse {
		res.BooksInUse = append(res.BooksInUse, &gen.BookInUse{
			BookInfo: &gen.Book{
				ID:     int32(bookInUse.BookInfo.ID),
				Title:  bookInUse.BookInfo.Title,
				Copies: int32(bookInUse.BookInfo.Copies),
				Author: bookInUse.BookInfo.Author,
				Issue:  timestamppb.New(bookInUse.BookInfo.Issue),
			},
			DateOfRent: timestamppb.New(bookInUse.DateOfRent),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetBooksInUseByReaderId(ctx context.Context, req *gen.ReaderIdRequest) (*gen.BooksInUseResponse, error) {
	booksInUse, err := s.useCase.GetBooksInUseByReaderId(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}

	res := &gen.BooksInUseResponse{}
	for _, bookInUse := range booksInUse {
		res.BooksInUse = append(res.BooksInUse, &gen.BookInUse{
			BookInfo: &gen.Book{
				ID:     int32(bookInUse.BookInfo.ID),
				Title:  bookInUse.BookInfo.Title,
				Copies: int32(bookInUse.BookInfo.Copies),
				Author: bookInUse.BookInfo.Author,
				Issue:  timestamppb.New(bookInUse.BookInfo.Issue),
			},
			DateOfRent: timestamppb.New(bookInUse.DateOfRent),
		})
	}

	return res, nil
}

func (s *GRPCServer) CountBookInUseByReaderId(ctx context.Context, req *gen.ReaderIdRequest) (*gen.CountResponse, error) {
	countBookInUse, err := s.useCase.CountBookInUseByReaderId(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}

	res := &gen.CountResponse{Count: int32(countBookInUse)}

	return res, nil
}

func (s *GRPCServer) CreateBookInUse(ctx context.Context, req *gen.BookInUseIdRequest) (*gen.MessageResponse, error) {

	err := s.useCase.CreateBookInUse(ctx, int(req.BookId), int(req.BookId))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book in use created"}
	return res, nil
}

func (s *GRPCServer) DeleteBookInUse(ctx context.Context, req *gen.BookInUseIdRequest) (*gen.MessageResponse, error) {

	err := s.useCase.DeleteBookInUse(ctx, int(req.BookId), int(req.BookId))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Book in use deleted"}
	return res, nil
}
