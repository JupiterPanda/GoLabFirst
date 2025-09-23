package booksinuse

import (
	"context"
	"fmt"
	"goproject/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository реализует доступ к данным по книгам
type Repository struct {
	db *pgxpool.Pool
}

// NewRepo создаёт новый репозиторий с пулом соединений к базе
func NewRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create добавляет в бд запись о новой книге у читателя
func (r *Repository) Create(ctx context.Context, bookInUse *models.BookInUse, readerId int) error {
	query := `INSERT INTO reader_books (book_id, reader_id, date_of_rent) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, bookInUse.BookInfo.ID, readerId, time.Now())
	return err
}

// GetAll возвращает все книги из базы
func (r *Repository) GetAll(ctx context.Context) ([]models.BookInUse, error) {
	rows, err := r.db.Query(ctx, `SELECT book_id, date_of_rent FROM reader_books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booksInUse []models.BookInUse

	for rows.Next() {
		var bookInUse models.BookInUse
		err = rows.Scan(&bookInUse.BookInfo, &bookInUse.DateOfRent)
		if err != nil {
			return nil, err
		}
		booksInUse = append(booksInUse, bookInUse)
	}
	return booksInUse, nil
}

// CountByReaderId ищет книги по ID клиента
func (r *Repository) CountByReaderId(ctx context.Context, readerId int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM reader_books WHERE reader_id = $1`
	err := r.db.QueryRow(ctx, query, readerId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetReadersIdsByBookId ищет клиентов, взявших книгу по ID книги
func (r *Repository) GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error) {
	rows, err := r.db.Query(ctx, `SELECT reader_id FROM reader_books WHERE book_id=$1`, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readerIds []int

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		readerIds = append(readerIds, id)
	}
	return readerIds, nil
}

// Delete удаляет из бд запись об аренде книги
func (r *Repository) Delete(ctx context.Context, readerId int, bookId int) error {
	query := `
        DELETE FROM reader_books 
        WHERE reader_id = $1 AND book_id = $2
    `
	cmdTag, err := r.db.Exec(ctx, query, readerId, bookId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no entry found to delete")
	}
	return nil
}

// GetBooksInUseByReaderId возвращает мап со значениями id книги: дата аренды
func (r *Repository) GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error) {
	rows, err := r.db.Query(ctx, `SELECT book_id, reader_id, date_of_rent FROM reader_books WHERE reader_id=$1`, readerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books map[int]time.Time
	var rentTime time.Time
	for rows.Next() {
		var book int
		err = rows.Scan(&book, &rentTime)
		if err != nil {
			return nil, err
		}
		books[book] = rentTime
	}
	return books, nil
}
