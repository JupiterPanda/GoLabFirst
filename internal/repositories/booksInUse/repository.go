package booksinuse

import (
	"context"
	"fmt"
	"goproject/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, bookInUse models.BookInUse, readerId int) error {

	query := `INSERT INTO reader_books (book_id, reader_id, date_of_rent) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, bookInUse.BookInfo.ID, readerId, time.Now())
	if err != nil {
		// TODO Проверка на вставку дубликата.
		return fmt.Errorf("[repo][Create] ошибка при запросе в БД: %w", err)
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]models.BookInUse, error) {
	rows, err := r.db.Query(ctx, `SELECT book_id, date_of_rent FROM reader_books`)
	if err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при запросе в БД: %w", err)
	}
	defer rows.Close()

	var booksInUse []models.BookInUse
	for rows.Next() {
		var bookInUse models.BookInUse
		var bookId int
		err = rows.Scan(&bookId, &bookInUse.DateOfRent)
		if err != nil {
			return nil, fmt.Errorf("[repo][GetAll] ошибка сканирования: %w", err)
		}
		bookInUse.BookInfo = models.Book{ID: bookId}
		booksInUse = append(booksInUse, bookInUse)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при проходе по rows: %w", err)
	}
	return booksInUse, nil
}

func (r *Repository) CountByReaderId(ctx context.Context, readerId int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM reader_books WHERE reader_id = $1`
	err := r.db.QueryRow(ctx, query, readerId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("[repo][CountByReaderId] ошибка при запросе в БД: %w", err)
	}
	return count, nil
}

func (r *Repository) GetReadersIdsByBookId(ctx context.Context, bookId int) ([]int, error) {
	rows, err := r.db.Query(ctx, `SELECT reader_id FROM reader_books WHERE book_id=$1`, bookId)
	if err != nil {
		return nil, fmt.Errorf("[repo][GetReadersIdsByBookId] ошибка при запросе в БД: %w", err)
	}
	defer rows.Close()

	var readerIds []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("[repo][GetReadersIdsByBookId] ошибка сканирования id: %w", err)
		}
		readerIds = append(readerIds, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo][GetReadersIdsByBookId] ошибка при проходе по rows: %w", err)
	}
	return readerIds, nil
}

func (r *Repository) Delete(ctx context.Context, readerId int, bookId int) error {
	query := `DELETE FROM reader_books WHERE reader_id = $1 AND book_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, readerId, bookId)
	if err != nil {
		return fmt.Errorf("[repo][Delete] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("[repo][Delete] запись об аренде книги не найдена")
	}
	return nil
}

func (r *Repository) GetBooksInUseByReaderId(ctx context.Context, readerId int) (map[int]time.Time, error) {
	rows, err := r.db.Query(ctx, `SELECT book_id, date_of_rent FROM reader_books WHERE reader_id=$1`, readerId)
	if err != nil {
		return nil, fmt.Errorf("[repo][GetBooksInUseByReaderId] ошибка при запросе в БД: %w", err)
	}
	defer rows.Close()

	books := make(map[int]time.Time)
	for rows.Next() {
		var bookId int
		var rentTime time.Time
		err = rows.Scan(&bookId, &rentTime)
		if err != nil {
			return nil, fmt.Errorf("[repo][GetBooksInUseByReaderId] ошибка сканирования: %w", err)
		}
		books[bookId] = rentTime
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo][GetBooksInUseByReaderId] ошибка при проходе по rows: %w", err)
	}
	return books, nil
}
