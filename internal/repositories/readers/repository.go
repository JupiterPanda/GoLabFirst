package readers

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrReaderNotFound = errors.New("читатель не найден")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, reader models.Reader) error {
	query := `INSERT INTO readers (name, number, address, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, reader.Name, reader.PhoneNumber, reader.Address, reader.DateOfBirth).Scan(&reader.ID)
	if err != nil {
		// TODO Проверка на вставку дубликата.
		return fmt.Errorf("[repo][Create] ошибка при запросе в БД: %w", err)
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Reader, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, number, address, date_of_birth FROM readers`)
	if err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при запросе в БД: %w", err)
	}
	defer rows.Close()

	var readers []models.Reader
	for rows.Next() {
		var reader models.Reader
		err = rows.Scan(&reader.ID, &reader.Name, &reader.PhoneNumber, &reader.Address, &reader.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("[repo][GetAll] ошибка сканирования данных читателя: %w", err)
		}
		readers = append(readers, reader)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo][GetAll] ошибка при проходе по rows: %w", err)
	}
	return readers, nil
}

func (r *Repository) GetIdByName(ctx context.Context, name string) (int, error) {
	var readerID int
	query := `SELECT id FROM readers WHERE name = $1`
	err := r.db.QueryRow(ctx, query, name).Scan(&readerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("[repo][GetIdByName] %w", ErrReaderNotFound)
		}
		return 0, fmt.Errorf("[repo][GetIdByName] ошибка при запросе в БД: %w", err)
	}
	return readerID, nil
}

func (r *Repository) Delete(ctx context.Context, reader models.Reader) error {
	query := `DELETE FROM readers WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, reader.ID)
	if err != nil {
		return fmt.Errorf("[repo][Delete] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("[repo][Delete] %w", ErrReaderNotFound)
	}
	return nil
}

func (r *Repository) UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error {
	query := `UPDATE readers SET number = $1, address = $2 WHERE id = $3`
	cmdTag, err := r.db.Exec(ctx, query, phoneNumber, address, readerId)
	if err != nil {
		return fmt.Errorf("[repo][UpdateContactInfo] ошибка при запросе в БД: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("[repo][UpdateContactInfo] читатель с ID %d: %w", readerId, ErrReaderNotFound)
	}
	return nil
}
