package handlers

import (
	"goproject/internal/models"
	"goproject/internal/services/books"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Service *books.Service
}

func NewBookHandler(service *books.Service) *BookHandler {
	return &BookHandler{Service: service}
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.Service.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var newBook models.Book
	err := c.Bind(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	err = h.Service.CreateBook(c.Request.Context(), &newBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create book"})
		return
	}
	c.JSON(http.StatusOK, "Success")
}

func (h *BookHandler) GetBookByTitle(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	book, err := h.Service.GetBookByTitle(c.Request.Context(), input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}
	c.JSON(http.StatusOK, book)
}
