package handlers

import (
	"goproject/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetReaderBooksSepGoodAndBad(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	okBooks, badBooks, err := h.useCase.GetReaderBooksSepGoodAndBad(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"okbooks":  okBooks,
		"badbooks": badBooks,
	})
}

func (h *Handler) RentBookByTitleAndReaderName(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	err := h.useCase.RentBookByTitleAndReaderName(c.Request.Context(), input.Name, input.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to rent book", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book rented"})
}

func (h *Handler) ReturnBookByTitleAndReaderName(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	err := h.useCase.ReturnBookByTitleAndReaderName(c.Request.Context(), input.Name, input.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to return book", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned"})
}

func (h *Handler) GetAllBooks(c *gin.Context) {
	books, err := h.useCase.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *Handler) GetBookByTitle(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	book, err := h.useCase.GetBookByTitle(c.Request.Context(), input.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *Handler) GetBookIdByTitle(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	id, err := h.useCase.GetBookIdByTitle(c.Request.Context(), input.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) GetBookByID(c *gin.Context) {
	var input struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	book, err := h.useCase.GetBookByID(c.Request.Context(), input.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *Handler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.CreateBook(c.Request.Context(), &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create book", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book created"})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.DeleteBook(c.Request.Context(), &book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete book", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func (h *Handler) CheckCopiesOfBookByID(c *gin.Context) {
	var input struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.CheckCopiesOfBookByID(c.Request.Context(), input.ID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "No copies", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Copies available"})
}

func (h *Handler) CheckCopiesOfBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.CheckCopiesOfBook(c.Request.Context(), &book)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "No copies", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Copies available"})
}

func (h *Handler) MinusCopyOfBookById(c *gin.Context) {
	var input struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.MinusCopyOfBookById(c.Request.Context(), input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to decrease copies", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Decreased copies"})
}

func (h *Handler) PlusCopyOfBookById(c *gin.Context) {
	var input struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.PlusCopyOfBookById(c.Request.Context(), input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to increase copies", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Increased copies"})
}

func (h *Handler) CreateBookInUse(c *gin.Context) {
	var input struct {
		BookID   int `json:"book_id" binding:"required"`
		ReaderId int `json:"reader_id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	var bookInUse models.BookInUse
	var err error
	bookPtr, err := h.useCase.GetBookByID(c.Request.Context(), input.BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create book in use", "error": err.Error()})
		return
	}
	bookInUse.BookInfo = *bookPtr
	bookInUse.DateOfRent = time.Now()

	err = h.useCase.CreateBookInUse(c.Request.Context(), &bookInUse, input.ReaderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create book in use", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book in use created"})
}

func (h *Handler) GetAllBooksInUse(c *gin.Context) {
	books, err := h.useCase.GetAllBooksInUse(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books in use", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books_in_use": books})
}

func (h *Handler) CountBookInUseByReaderId(c *gin.Context) {
	var input struct {
		ReaderId int `json:"reader_id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	count, err := h.useCase.CountBookInUseByReaderId(c.Request.Context(), input.ReaderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to count books in use", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *Handler) GetReadersIdsByBookId(c *gin.Context) {
	var input struct {
		BookId int `json:"book_id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	ids, err := h.useCase.GetReadersIdsByBookId(c.Request.Context(), input.BookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get reader ids", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reader_ids": ids})
}

func (h *Handler) GetBooksInUseByReaderId(c *gin.Context) {
	var input struct {
		ReaderId int `json:"reader_id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	books, err := h.useCase.GetBooksInUseByReaderId(c.Request.Context(), input.ReaderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get books in use", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books_in_use": books})
}

func (h *Handler) DeleteBookInUse(c *gin.Context) {
	var input struct {
		ReaderId int `json:"reader_id" binding:"required"`
		BookId   int `json:"book_id" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.DeleteBookInUse(c.Request.Context(), input.ReaderId, input.BookId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete book in use", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book in use deleted"})
}

func (h *Handler) GetAllReaders(c *gin.Context) {
	readers, err := h.useCase.GetAllReaders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get readers", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"readers": readers})
}

func (h *Handler) CreateReader(c *gin.Context) {
	var reader models.Reader
	if err := c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.CreateReader(c.Request.Context(), &reader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create reader", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reader created"})
}

func (h *Handler) GetReaderIdByName(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	id, err := h.useCase.GetReaderIdByName(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Reader not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) DeleteReader(c *gin.Context) {
	var reader models.Reader
	if err := c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.DeleteReader(c.Request.Context(), &reader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete reader", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reader deleted"})
}

func (h *Handler) UpdateReaderContactInfo(c *gin.Context) {
	var input struct {
		ReaderId    int    `json:"reader_id" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Address     string `json:"address" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}
	err := h.useCase.UpdateReaderContactInfo(c.Request.Context(), input.ReaderId, input.PhoneNumber, input.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update contact info", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contact info updated"})
}
