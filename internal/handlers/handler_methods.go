package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetReaderBooksSepGoodAndBad(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	okBooks, badBooks, err := h.useCase.GetReaderBooksSepGoodAndBad(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"okbooks":  okBooks,
		"badbooks": badBooks,
	})
}

func (h *Handler) RentBookByTitleAndReaderName(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := h.useCase.RentBookByTitleAndReaderName(c.Request.Context(), input.Name, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book rented"})
}

func (h *Handler) ReturnBookByTitleAndReaderName(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := h.useCase.ReturnBookByTitleAndReaderName(c.Request.Context(), input.Name, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book returned"})
}

func (h *Handler) GetAllBooks(c *gin.Context) {

	books, err := h.useCase.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByTitle(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	books, err := h.useCase.GetBookByTitle(c.Request.Context(), input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get books"})
		return
	}

	c.JSON(http.StatusOK, books)
}
