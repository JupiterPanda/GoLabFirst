package app

import (
	"goproject/internal/handlers"

	"github.com/gin-gonic/gin"
)

func initRouter(handler *handlers.Handler) *gin.Engine {
	router := gin.Default()

	// Чтение, аренда, возврат
	router.GET("/reader/books", handler.GetReaderBooksSepGoodAndBad)
	router.PATCH("/rent", handler.RentBookByTitleAndReaderName)
	router.PATCH("/return", handler.ReturnBookByTitleAndReaderName)

	// Книги CRUD
	router.GET("/books", handler.GetAllBooks)
	router.GET("/book/title", handler.GetBookByTitle)
	router.GET("/book/id/title", handler.GetBookIdByTitle)
	router.GET("/book/id", handler.GetBookByID)
	router.POST("/book", handler.CreateBook)
	router.DELETE("/book", handler.DeleteBook)

	// Копии книги
	router.POST("/book/check/id", handler.CheckCopiesOfBookByID)
	router.POST("/book/check", handler.CheckCopiesOfBook)
	router.PATCH("/book/minus", handler.SubtractCopyOfBookById)
	router.PATCH("/book/plus", handler.AddCopyOfBookById)

	// Операции с InUse
	router.POST("/reader/book", handler.CreateBookInUse)
	router.GET("/reader/book", handler.GetAllBooksInUse)
	router.GET("/reader/book/count", handler.CountBookInUseByReaderId)
	router.GET("/reader/book/id", handler.GetBooksInUseByReaderId)
	router.DELETE("/reader/book", handler.DeleteBookInUse)

	// CRUD для читателей
	router.GET("/readers", handler.GetAllReaders)
	router.GET("/reader/id", handler.GetReaderIdByName)
	router.POST("/reader", handler.CreateReader)
	router.DELETE("/reader", handler.DeleteReader)
	router.PATCH("/reader/contact", handler.UpdateReaderContactInfo)
	router.GET("/reader/id/book", handler.GetReadersIdsByBookId)

	return router
}
