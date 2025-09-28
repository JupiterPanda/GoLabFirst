package readers

import (
	"context"
	"goproject/internal/models"
	constants "goproject/internal/package"
	repository "goproject/internal/repository/repos"
	"log"
	"time"
)

type Service struct {
	repo *repository.ReaderRepository
}

// Конструктор сервиса читателя
func NewService(repo *repository.ReaderRepository) *Service {
	return &Service{repo: repo}
}

// Получить все книги
func (s *Service) GetAllReaders(ctx context.Context) ([]models.Reader, error) {
	return s.repo.GetAll(ctx)
}

// Добавить читателя
func (s *Service) CreateReader(ctx context.Context, reader *models.Reader) error {
	return s.repo.Create(ctx, reader)
}

// Получить все книги у пользователя
func (s *Service) GetReaderBooks(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error) {
	readerId, err := s.repo.GetReaderIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	booksInUse, err := s.repo.GetReaderBooksByID(ctx, readerId)
	if err != nil {
		log.Printf("unknown error")
	}

	var okbooks, badbooks []models.BookInUse

	for _, book := range booksInUse {
		if time.Since(book.DateOfRent) <= constants.TimeToExpire {
			okbooks = append(okbooks, book)
		} else {
			badbooks = append(badbooks, book)
		}
	}
	return okbooks, badbooks, nil
}

/*
func (s *Service, r *Service) RentBookByTitle(ctx context.Context, name, title string) error {

	bookId, err := GetByTitle(ctx, title)
	if err != nil {
		log.Printf("book not found")
	}
	ok := repository.CheckCopiesByID(ctx, bookId)

	readerId, err := s.repo.GetReaderIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	CheckNumOfBooksOfReader(ctx, readerId)

	if name == "" || title == "" {
		return errors.New("name and title must be provided")
	}
	return s.repo.RentBook(ctx, readerId, bookId)
} */

/* func (s *Service) ReturnBookByTitle(ctx context.Context, name, title string) error {
	if name == "" || title == "" {
		return errors.New("name and title must be provided")
	}
	return s.repo.ReturnBook(ctx, name, title)
} */
