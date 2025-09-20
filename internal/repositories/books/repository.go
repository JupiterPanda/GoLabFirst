package books

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository реализует доступ к данным по книгам
type Repository struct {
	db *pgxpool.Pool
}

// New создаёт новый репозиторий с пулом соединений к базе
func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create добавляет новую книгу в базу
func (r *Repository) Create(ctx context.Context, book *models.Book) error {
	query := `INSERT INTO books (title, author, issue, copies) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.Issue, book.Copies).Scan(&book.ID)
	return err
}

// GetAll возвращает все книги из базы
func (r *Repository) GetAll(ctx context.Context) ([]models.Book, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, author, issue, copies FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetByTitle ищет книгу по названию
func (r *Repository) GetByTitle(ctx context.Context, title string) (*models.Book, error) {
	var book models.Book
	query := `SELECT id, title, author, issue, copies FROM books WHERE title=$1`
	err := r.db.QueryRow(ctx, query, title).Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

// CheckCopiesByID ищет книгу по названию
func (r *Repository) CheckCopiesByID(ctx context.Context, id int) error {
	// Найти книгу по названию с достаточным количеством копий
	var copies int
	err := r.db.QueryRow(ctx, "SELECT id, copies FROM books WHERE id=$1", id).Scan(&copies)
	if err != nil {
		return errors.New("book not found")
	}
	if copies <= 0 {
		return errors.New("this book are out of stock")
	}
	return nil
}

// Delete удаляет из бд книгу (!!! Удалит книги и в таблице reader_books!!!)
func (r *Repository) Delete(ctx context.Context, book *models.Book) error {
	query := `
        DELETE FROM books 
        WHERE title = $1
    `
	cmdTag, err := r.db.Exec(ctx, query, book.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no entry found to delete")
	}
	return nil
}

// PlusCopyById прибавляет кол-во свободных копий книге
func (r *Repository) PlusCopyById(ctx context.Context, book *models.Book) error {
	query := `UPDATE books SET copies = copies + 1 WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, book.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("book with id %d not found", book.ID)
	}
	return nil
}

// MinusCopyById убавляет кол-во свободных копий книги
func (r *Repository) MinusCopyById(ctx context.Context, book *models.Book) error {
	query := `UPDATE books SET copies = copies - 1 WHERE id = $1 AND copies > 0`
	cmdTag, err := r.db.Exec(ctx, query, book.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("book with id %d not found or no available copies", book.ID)
	}
	return nil
}
