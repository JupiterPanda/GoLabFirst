package repository

import (
	"context"
	"errors"
	"goproject/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BookRepository реализует доступ к данным по книгам
type BookRepository struct {
	db *pgxpool.Pool
}

// NewBookRepository создаёт новый репозиторий с пулом соединений к базе
func NewBookRepository(db *pgxpool.Pool) *BookRepository {
	return &BookRepository{db: db}
}

// Create добавляет новую книгу в базу
func (r *BookRepository) Create(ctx context.Context, book *models.Book) error {
	query := `INSERT INTO books (title, author, issue, copies) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.Issue, book.Copies).Scan(&book.ID)
	return err
}

// GetAll возвращает все книги из базы
func (r *BookRepository) GetAll(ctx context.Context) ([]models.Book, error) {
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
func (r *BookRepository) GetByTitle(ctx context.Context, title string) (*models.Book, error) {
	var book models.Book
	query := `SELECT id, title, author, issue, copies FROM books WHERE title=$1`
	err := r.db.QueryRow(ctx, query, title).Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

// GetByTitle ищет книгу по названию
func (r *BookRepository) CheckCopiesByID(ctx context.Context, id int) error {
	// Найти книгу по названию с достаточным количеством копий
	var copies int
	err := r.db.QueryRow(ctx, "SELECT id, copies FROM book WHERE id=$1", id).Scan(&copies)
	if err != nil {
		return errors.New("book not found")
	}
	if copies <= 0 {
		return errors.New("this book are out of stock")
	}
	return nil
}
