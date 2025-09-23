package readers

import (
	"context"
	"goproject/internal/models"
	repository "goproject/internal/repositories/readers"
)

type Service struct {
	repo *repository.Repository
}

// NewService Конструктор сервиса читателя
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// GetAllReaders Получить все книги
func (s *Service) GetAll(ctx context.Context) ([]models.Reader, error) {
	return s.repo.GetAll(ctx)
}

// CreateReader Добавить читателя
func (s *Service) Create(ctx context.Context, reader *models.Reader) error {
	return s.repo.Create(ctx, reader)
}

// GetIdByName Метод для получения ID читателя по имени
func (s *Service) GetIdByName(ctx context.Context, name string) (int, error) {
	return s.repo.GetIdByName(ctx, name)
}

// Delete уладить читателя из readers
func (s *Service) Delete(ctx context.Context, readerId int) error {
	return s.repo.Delete(ctx, readerId)
}

// UpdateContactInfo обновит значения номера телефона или адреса
func (s *Service) UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error {
	return s.repo.UpdateContactInfo(ctx, readerId, phoneNumber, address)
}

// GetReaderBooks Получить все книги у пользователя
/* func (s *Service) GetReaderBooks(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error) {
	readerId, err := s.repo.GetReaderIdByName(ctx, name)
	if err != nil {
		log.Printf("reader not found")
	}

	booksInUse, err := s.repo.GetReaderBooksByID(ctx, readerId)
	if err != nil {
		log.Printf("unknown error")
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
} */

/*func (s *Service, r *Service) RentBookByTitle(ctx context.Context, name, title string) error {

	bookId, err := GetByTitle(ctx, title)
	if err != nil {
		log.Printf("book not found")
	}
	ok := repositories.CheckCopiesByID(ctx, bookId)

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
