package libgrpc

import (
	"context"
	"goproject/protos/gen/librarypb"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) RentBookByTitleAndReaderName(ctx context.Context, req *librarypb.RentOrReturnBookRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.RentBookByTitleAndReaderName(ctx, req.Name, req.Title)
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: "Book rented successfully"}
	return res, nil
}

func (s *GRPCServer) ReturnBookByTitleAndReaderName(ctx context.Context, req *librarypb.RentOrReturnBookRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.ReturnBookByTitleAndReaderName(ctx, req.Name, req.Title)
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: "Book returned successfully"}
	return res, nil
}

func (s *GRPCServer) GetReaderBooksSepGoodAndBad(ctx context.Context, req *librarypb.ReaderNameRequest) (*librarypb.ReadersBooksResponse, error) {
	okBooks, badBooks, err := s.useCase.GetReaderBooksSepGoodAndBad(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	res := &librarypb.ReadersBooksResponse{}

	for _, bookInUse := range okBooks {
		res.OkBooks = append(res.OkBooks, &librarypb.BookInUse{
			BookInfo: &librarypb.Book{
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
		res.BadBooks = append(res.BadBooks, &librarypb.BookInUse{
			BookInfo: &librarypb.Book{
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

func (s *GRPCServer) GetReadersIdsByBookId(ctx context.Context, req *librarypb.BookIdRequest) (*librarypb.ReaderIdsResponse, error) {
	listOfBooksInUseOfReader, err := s.useCase.GetReadersIdsByBookId(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	int32ListOfBooksInUseOfReader := make([]int32, len(listOfBooksInUseOfReader), 3)
	for i, id := range listOfBooksInUseOfReader {
		int32ListOfBooksInUseOfReader[i] = int32(id)
	}
	res := &librarypb.ReaderIdsResponse{ReaderIds: int32ListOfBooksInUseOfReader}
	return res, nil
}

func (s *GRPCServer) GetAllBooksInUse(ctx context.Context, _ *emptypb.Empty) (*librarypb.BooksInUseResponse, error) {
	booksInUse, err := s.useCase.GetAllBooksInUse(ctx)
	if err != nil {
		return nil, err
	}

	res := &librarypb.BooksInUseResponse{}
	for _, bookInUse := range booksInUse {
		res.BooksInUse = append(res.BooksInUse, &librarypb.BookInUse{
			BookInfo: &librarypb.Book{
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

func (s *GRPCServer) GetBooksInUseByReaderId(ctx context.Context, req *librarypb.ReaderIdRequest) (*librarypb.BooksInUseResponse, error) {
	booksInUse, err := s.useCase.GetBooksInUseByReaderId(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}

	res := &librarypb.BooksInUseResponse{}
	for _, bookInUse := range booksInUse {
		res.BooksInUse = append(res.BooksInUse, &librarypb.BookInUse{
			BookInfo: &librarypb.Book{
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

func (s *GRPCServer) CountBookInUseByReaderId(ctx context.Context, req *librarypb.ReaderIdRequest) (*librarypb.CountResponse, error) {
	countBookInUse, err := s.useCase.CountBookInUseByReaderId(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}

	res := &librarypb.CountResponse{Count: int32(countBookInUse)}

	return res, nil
}

func (s *GRPCServer) CreateBookInUse(ctx context.Context, req *librarypb.BookInUseIdRequest) (*librarypb.MessageResponse, error) {

	err := s.useCase.CreateBookInUse(ctx, int(req.BookId), int(req.BookId))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: "Book in use created"}
	return res, nil
}

func (s *GRPCServer) DeleteBookInUse(ctx context.Context, req *librarypb.BookInUseIdRequest) (*librarypb.MessageResponse, error) {

	err := s.useCase.DeleteBookInUse(ctx, int(req.BookId), int(req.BookId))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: "Book in use deleted"}
	return res, nil
}
