package repository

import (
	"context"
	"errors"
	"goproject/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReaderRepository реализует доступ к данным по читателям
type ReaderRepository struct {
	db *pgxpool.Pool
}

// Конструктор для создания нового репозитория читателя
func NewReaderRepository(db *pgxpool.Pool) *ReaderRepository {
	return &ReaderRepository{db: db}
}

// Create добавляет нового читателя в базу
func (r *ReaderRepository) Create(ctx context.Context, reader *models.Reader) error {
	query := `INSERT INTO readers (name, number, address, dateofbirth) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, reader.Name, reader.PhoneNumber, reader.Address, reader.DateOfBirth).Scan(&reader.ID)
}

// GetAll возвращает всех читателей из базы
func (r *ReaderRepository) GetAll(ctx context.Context) ([]models.Reader, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, number, address, dateofbirth FROM readers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readers []models.Reader
	for rows.Next() {
		var reader models.Reader
		err = rows.Scan(&reader.ID, &reader.Name, &reader.PhoneNumber, &reader.Address, &reader.DateOfBirth)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

// Метод для получения ID читателя по имени (можно перенести в GetReaderBooksByID)
func (r *ReaderRepository) GetReaderIdByName(ctx context.Context, name string) (int, error) {
	var readerID int
	err := r.db.QueryRow(ctx, "SELECT id FROM readers WHERE name = $1", name).Scan(&readerID)
	if err != nil {
		return 0, errors.New("reader not found")
	}
	return readerID, nil
}

// Метод для получения книг по ID читателя
func (r *ReaderRepository) GetReaderBooksByID(ctx context.Context, readerID int) ([]models.BookInUse, error) {

	// Получаем книги из связующей таблицы readerBooks с join к таблице book
	query := `
        SELECT b.id, b.title, b.author, b.issue, b.copies, rb.dateofrent
        FROM readerBooks rb
        JOIN books b ON rb.book_id = b.id
        WHERE rb.reader_id = $1
    `
	rows, err := r.db.Query(ctx, query, readerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booksInUse []models.BookInUse
	for rows.Next() {
		var book models.Book
		var dateOfRent time.Time
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies, &dateOfRent)
		if err != nil {
			return nil, err
		}
		booksInUse = append(booksInUse, models.BookInUse{
			NameOfBook: book,
			DateOfRent: dateOfRent,
		})
	}
	return booksInUse, nil
}

// GetByTitle ищет книгу по названию
func (r *BookRepository) CheckNumOfBooksOfReader(ctx context.Context, readerID int) error {
	// Найти книгу по названию с достаточным количеством копий
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM readerBook WHERE reader_id=$1", readerID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= 3 {
		return errors.New("too many books rented")
	}
	return nil
}

// Метод для аренды книги по ID читателя и книги
func (r *ReaderRepository) RentBook(ctx context.Context, readerID, bookID int) error {

	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Уменьшаем количество доступных копий книги
	_, err = tx.Exec(ctx, "UPDATE book SET copies = copies - 1 WHERE id = $1", bookID)
	if err != nil {
		return err
	}

	// Добавляем запись о аренде в связующую таблицу
	_, err = tx.Exec(ctx,
		"INSERT INTO readerBook (book_id, reader_id, dateofrent) VALUES ($1, $2, $3)",
		bookID, readerID, time.Now())
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	return tx.Commit(ctx)
}
