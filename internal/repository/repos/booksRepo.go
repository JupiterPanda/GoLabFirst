package repository

import (
	"context"
	"errors"
	"goproject/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
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
	query := `INSERT INTO book (title, author, issue, copies) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.Issue, book.Copies).Scan(&book.ID)
	return err
}

// GetAll возвращает все книги из базы
func (r *BookRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, author, issue, copies FROM book`)
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
	query := `SELECT id, title, author, issue, copies FROM book WHERE title=$1`
	err := r.db.QueryRow(ctx, query, title).Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return &book, nil
}
