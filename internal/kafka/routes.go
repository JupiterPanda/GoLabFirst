package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"goproject/internal/models"
)

type handlerFunc func(ctx context.Context, payload json.RawMessage) error

func (c *Consumer) routes() map[string]handlerFunc {
	return map[string]handlerFunc{
		"rent_book":   c.rentBook,
		"return_book": c.returnBook,
		"create_book": c.createBook,
		"delete_book": c.deleteBook,
		"add_copy":    c.addCopy,
		"minus_copy":  c.subtractCopy,
	}
}

type rentPayload struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (c *Consumer) rentBook(ctx context.Context, payload json.RawMessage) error {
	var p rentPayload
	if err := decode(payload, &p); err != nil {
		return err
	}
	if p.Name == "" || p.Title == "" {
		return fmt.Errorf("rent_book: name and title are required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// ровно тот же вызов, что делает gRPC-хендлер
	return c.useCase.RentBookByTitleAndReaderName(ctx, p.Name, p.Title)
}

func (c *Consumer) returnBook(ctx context.Context, payload json.RawMessage) error {
	var p rentPayload
	if err := decode(payload, &p); err != nil {
		return err
	}
	if p.Name == "" || p.Title == "" {
		return fmt.Errorf("return_book: name and title are required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.useCase.ReturnBookByTitleAndReaderName(ctx, p.Name, p.Title)
}

func (c *Consumer) createBook(ctx context.Context, payload json.RawMessage) error {
	var book models.Book // у models.Book уже есть json-теги, отдельный DTO не нужен
	if err := decode(payload, &book); err != nil {
		return err
	}
	if book.Title == "" {
		return fmt.Errorf("create_book: title is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.useCase.CreateBook(ctx, book)
}

type idPayload struct {
	ID int `json:"id"`
}

func (c *Consumer) deleteBook(ctx context.Context, payload json.RawMessage) error {
	var p idPayload
	if err := decode(payload, &p); err != nil {
		return err
	}
	return c.useCase.DeleteBook(ctx, p.ID)
}

func (c *Consumer) addCopy(ctx context.Context, payload json.RawMessage) error {
	var p idPayload
	if err := decode(payload, &p); err != nil {
		return err
	}
	return c.useCase.AddCopyOfBookById(ctx, p.ID)
}

func (c *Consumer) subtractCopy(ctx context.Context, payload json.RawMessage) error {
	var p idPayload
	if err := decode(payload, &p); err != nil {
		return err
	}
	return c.useCase.SubtractCopyOfBookById(ctx, p.ID)
}

// decode — строгая десериализация: опечатки в именах полей ловим сразу, а не молча
func decode(payload json.RawMessage, dst any) error {
	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	return nil
}
