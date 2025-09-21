package readers

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository реализует доступ к данным по читателям
type Repository struct {
	db *pgxpool.Pool
}

// NewRepo Конструктор для создания нового репозитория читателя
func NewRepo(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create добавляет нового читателя в базу
func (r *Repository) Create(ctx context.Context, reader *models.Reader) error {
	query := `INSERT INTO readers (name, number, address, date_of_birth) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, reader.Name, reader.PhoneNumber, reader.Address, reader.DateOfBirth).Scan(&reader.ID)
}

// GetAll возвращает всех читателей из базы
func (r *Repository) GetAll(ctx context.Context) ([]models.Reader, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, number, address, date_of_birth FROM readers`)
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

// GetIdByName Метод для получения ID читателя по имени
func (r *Repository) GetIdByName(ctx context.Context, name string) (int, error) {
	var readerID int
	err := r.db.QueryRow(ctx, "SELECT id FROM readers WHERE name = $1", name).Scan(&readerID)
	if err != nil {
		return 0, errors.New("reader not found")
	}
	return readerID, nil
}

// Delete Удалит читателя из readers
func (r *Repository) Delete(ctx context.Context, readerId int) error {
	query := `DELETE FROM readers WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, readerId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reader with id %d not found", readerId)
	}
	return nil
}

// UpdateContactInfo обновит значения номера телефона или адреса
func (r *Repository) UpdateContactInfo(ctx context.Context, readerId int, phoneNumber string, address string) error {
	query := `
        UPDATE readers
        SET number = $1, address = $2
        WHERE id = $3
    `
	cmdTag, err := r.db.Exec(ctx, query, phoneNumber, address, readerId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reader with id %d not found", readerId)
	}
	return nil
}
