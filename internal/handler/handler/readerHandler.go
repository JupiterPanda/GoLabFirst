package handlers

import (
	"goproject/internal/models"
	"goproject/internal/services/readers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReaderHandler struct {
	Service *readers.Service
}

func NewReaderHandler(service *readers.Service) *ReaderHandler {
	return &ReaderHandler{Service: service}
}

// Получить всех читателей
func (h *ReaderHandler) GetAllReaders(c *gin.Context) {
	readers, err := h.Service.GetAllReaders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get readers"})
		return
	}
	c.JSON(http.StatusOK, readers)
}

// Добавить нового читателя
func (h *ReaderHandler) CreateReader(c *gin.Context) {
	var newReader models.Reader
	if err := c.BindJSON(&newReader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	err := h.Service.CreateReader(c.Request.Context(), &newReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create reader"})
		return
	}
	c.JSON(http.StatusCreated, newReader)
}

// Получить книги, взятые читателем с разделением на просроченные и нет
func (h *ReaderHandler) GetReaderBooks(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	okbooks, badbooks, err := h.Service.GetReaderBooks(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Reader not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"okbooks":  okbooks,
		"badbooks": badbooks,
	})
}
