package handlers

import (
	"context"
	"goproject/internal/models"
)

type Handler struct {
	useCase UseCase
}

type UseCase interface {
	GetReaderBooksSepGoodAndBad(ctx context.Context, name string) ([]models.BookInUse, []models.BookInUse, error)
	RentBookByTitleAndReaderName(ctx context.Context, name, title string) error
	ReturnBookByTitleAndReaderName(ctx context.Context, name, title string) error
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}
