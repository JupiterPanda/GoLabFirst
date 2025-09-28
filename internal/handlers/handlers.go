package handlers

import (
	"goproject/internal/services/books"
	"goproject/internal/services/readers"
)

type Handlers struct {
	bookService   books.Service
	readerService readers.Service
}

func NewHandlers(
	bookService books.Service,
	readerService readers.Service,
) *Handlers {
	return &Handlers{
		bookService:   bookService,
		readerService: readerService,
	}
}

// func (h *Handlers) Init(api *gin.RouterGroup) {
//     v1 := api.Group("/v1")
//     {
//         h.initStudentsRoutes(v1)
//     }
// }
