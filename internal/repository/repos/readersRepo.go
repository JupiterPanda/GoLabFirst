package repository

import (
	"context"
	"errors"
	"goproject/internal/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ReaderRepository реализует доступ к данным по читателям
type ReaderRepository struct {
	db *pgxpool.Pool
}

func NewReaderRepository(db *pgxpool.Pool) *ReaderRepository {
	return &ReaderRepository{db: db}
}

// Create добавляет нового читателя в базу
func (r *ReaderRepository) Create(ctx context.Context, reader *models.Reader) error {
	query := `INSERT INTO reader (name, number, address, dateofbirth) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, reader.Name, reader.PhoneNumber, reader.Address, reader.DateOfBirth).Scan(&reader.ID)
}

// GetAll возвращает всех читателей из базы
func (r *ReaderRepository) GetAll(ctx context.Context) ([]models.Reader, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, number, address, dateofbirth FROM reader`)
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

// Метод для получения книг, взятых читателем по имени
func (r *ReaderRepository) GetBooksInUseByName(ctx context.Context, name string) ([]models.BookInUse, error) {
	// 1. Сначала находим ID читателя по имени
	var readerID int
	err := r.db.QueryRow(ctx, "SELECT id FROM reader WHERE name = $1", name).Scan(&readerID)
	if err != nil {
		return nil, errors.New("reader not found")
	}

	// 2. Получаем книги из связующей таблицы readerBook с join к таблице book
	query := `
        SELECT b.id, b.title, b.author, b.issue, b.copies, rb.dateofrent
        FROM readerBook rb
        JOIN book b ON rb.book_id = b.id
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

// Метод для оформления аренды книги по имени читателя и названию книги
func (r *ReaderRepository) RentBook(ctx context.Context, readerName, bookTitle string) error {
	// Найти книгу по названию с достаточным количеством копий
	var bookID, copies int
	err := r.db.QueryRow(ctx, "SELECT id, copies FROM book WHERE title=$1", bookTitle).Scan(&bookID, &copies)
	if err != nil {
		return errors.New("book not found")
	}
	if copies <= 0 {
		return errors.New("all books are rented")
	}

	// Найти читателя по имени
	var readerID int
	err = r.db.QueryRow(ctx, "SELECT id FROM reader WHERE name=$1", readerName).Scan(&readerID)
	if err != nil {
		return errors.New("reader not found")
	}

	// Проверить сколько книг сейчас взято читателем
	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM readerBook WHERE reader_id=$1", readerID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= 3 {
		return errors.New("too many books rented")
	}

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
