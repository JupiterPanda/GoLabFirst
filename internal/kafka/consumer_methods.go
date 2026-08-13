package kafka

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/models"
	"goproject/protos/gen"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

func (c *Consumer) handle(ctx context.Context, rec *kgo.Record) error {
	var event gen.Event
	if err := proto.Unmarshal(rec.Value, &event); err != nil {
		return fmt.Errorf("invalid protobuf event: %w", err)
	}
	log.Printf("kafka: received event (event=%s) with payload: %s", event.EventId, event.GetPayload())
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Роутинг по oneof: распаковка и вызов метода UseCase.
	switch p := event.GetPayload().(type) {

	case *gen.Event_RentBook:
		req := p.RentBook
		if req.GetName() == "" || req.GetTitle() == "" {
			return errors.New("rent_book: name and title are required")
		}
		return c.useCase.RentBookByTitleAndReaderName(ctx, req.GetName(), req.GetTitle())

	case *gen.Event_ReturnBook:
		req := p.ReturnBook
		if req.GetName() == "" || req.GetTitle() == "" {
			return errors.New("return_book: name and title are required")
		}
		return c.useCase.ReturnBookByTitleAndReaderName(ctx, req.GetName(), req.GetTitle())

	case *gen.Event_CreateBook:
		book, err := bookFromProto(p.CreateBook.GetBook())
		if err != nil {
			return fmt.Errorf("create_book: %w", err)
		}
		return c.useCase.CreateBook(ctx, book)

	case *gen.Event_DeleteBook:
		return c.useCase.DeleteBook(ctx, int(p.DeleteBook.GetId()))

	case *gen.Event_AddCopy:
		return c.useCase.AddCopyOfBookById(ctx, int(p.AddCopy.GetId()))

	case *gen.Event_MinusCopy:
		return c.useCase.SubtractCopyOfBookById(ctx, int(p.MinusCopy.GetId()))

	case nil:
		return errors.New("event has empty payload")

	default:
		return fmt.Errorf("unsupported payload type %T", p)
	}
}

func bookFromProto(b *gen.Book) (models.Book, error) {
	if b == nil {
		return models.Book{}, errors.New("book is nil")
	}
	if b.GetTitle() == "" {
		return models.Book{}, errors.New("title is required")
	}

	book := models.Book{
		ID:     int(b.GetID()),
		Title:  b.GetTitle(),
		Copies: int(b.GetCopies()),
		Author: b.GetAuthor(),
	}
	if ts := b.GetIssue(); ts != nil {
		book.Issue = ts.AsTime() // у nil AsTime() дал бы 1970 год
	}

	return book, nil
}
