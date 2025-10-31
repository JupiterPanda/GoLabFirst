package books

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"
	constants "goproject/internal/package"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, book models.Book) error {
	query := `INSERT INTO books (title, author, issue, copies) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, book.Title, book.Author, book.Issue, book.Copies)
	if err != nil {
		// TODO Проверка на вставку дубликата.
		return fmt.Errorf("[repo][Create] ошибка при запросе в БД: %w", err)
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Book, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, author, issue, copies FROM books`)
	if err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при запросе в БД: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
		if err != nil {
			return nil, fmt.Errorf("[repo][GetAll] ошибка сканирования данных книги: %w", err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при проходе по rows: %w", err)
	}
	return books, nil
}

func (r *Repository) GetByTitle(ctx context.Context, title string) (models.Book, error) {
	var book models.Book
	query := `SELECT id, title, author, issue, copies FROM books WHERE title=$1`
	err := r.db.QueryRow(ctx, query, title).Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, fmt.Errorf("[repo][GetByTitle] %w", constants.ErrBookNotFound)
		}
		return models.Book{}, fmt.Errorf("[repo][GetByTitle] ошибка при запросе в БД: %w", err)
	}
	return book, nil
}

func (r *Repository) GetIdByTitle(ctx context.Context, title string) (int, error) {
	var id int
	query := `SELECT id FROM books WHERE title=$1`
	err := r.db.QueryRow(ctx, query, title).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("[repo][GetIdByTitle] %w", constants.ErrBookNotFound)
		}
		return 0, fmt.Errorf("[repo][GetIdByTitle] ошибка при запросе в БД: %w", err)
	}
	return id, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (models.Book, error) {
	var book models.Book
	query := `SELECT id, title, author, issue, copies FROM books WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Issue, &book.Copies)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, fmt.Errorf("[repo][GetByID] %w", constants.ErrBookNotFound)
		}
		return models.Book{}, fmt.Errorf("[repo][GetByID] ошибка при запросе в БД: %w", err)
	}
	return book, nil
}

func (r *Repository) CheckCopiesByID(ctx context.Context, id int) error {
	var copies int
	query := `SELECT copies FROM books WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(&copies)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("[repo][CheckCopiesByID] %w", constants.ErrBookNotFound)
		}
		return fmt.Errorf("[repo][CheckCopiesByID] ошибка при запросе в БД: %w", err)
	}
	if copies <= 0 {
		return fmt.Errorf("[repo][CheckCopiesByID] %w", constants.ErrBookOutOfStock)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("[repo][Delete] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("[repo][Delete] %w", constants.ErrBookNotFound)
	}
	return nil
}

func (r *Repository) AddCopyById(ctx context.Context, id int) error {
	query := `UPDATE books SET copies = copies + 1 WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("[repo][AddCopyById] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("[repo][AddCopyById] %w", constants.ErrBookNotFound)
	}
	return nil
}

func (r *Repository) SubtractCopyById(ctx context.Context, id int) error {
	query := `UPDATE books SET copies = copies - 1 WHERE id = $1 AND copies > 0`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("[repo][SubtractCopyById] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		if err := r.CheckCopiesByID(ctx, id); err != nil {
			if errors.Is(err, constants.ErrBookNotFound) {
				return fmt.Errorf("[repo][SubtractCopyById] %w", constants.ErrBookNotFound)
			}
			if errors.Is(err, constants.ErrBookOutOfStock) {
				return fmt.Errorf("[repo][SubtractCopyById] %w", constants.ErrBookOutOfStock)
			}
		}
		return fmt.Errorf("[repo][SubtractCopyById] не удалось уменьшить копии")
	}
	return nil
}
