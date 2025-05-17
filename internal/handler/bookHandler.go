package handler

import (
	"github.com/gin-gonic/gin"
	"library-Mongo/internal/domain"
	"library-Mongo/internal/usecase"
	"library-Mongo/internal/usecase/dto"
	"net/http"
)

type BookHandler struct {
	bookUC usecase.BookUC
}

func NewBookHandler(bookUC usecase.BookUC) *BookHandler {
	return &BookHandler{bookUC: bookUC}
}

// CreateBook godoc
// @Summary Добавить новую книгу
// @Tags books
// @Accept json
// @Produce json
// @Param input body dto.CreateBookInput true "Данные книги"
// @Success 200 {object} domain.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var input dto.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	book, err := h.bookUC.CreateBook(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// UpdateBook godoc
// @Summary Обновить книгу
// @Tags books
// @Accept json
// @Produce json
// @Param input body dto.UpdateBookInput true "Обновляемые поля"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	var input dto.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	if err := h.bookUC.UpdateBook(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteBook godoc
// @Summary Удалить книгу
// @Tags books
// @Produce json
// @Param id path string true "ID книги"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := h.bookUC.DeleteBook(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// GetBookByID godoc
// @Summary Получить книгу по ID
// @Tags books
// @Produce json
// @Param id path string true "ID книги"
// @Success 200 {object} domain.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	book, err := h.bookUC.GetBookByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// SearchBooks godoc
// @Summary Поиск книг
// @Tags books
// @Produce json
// @Param title query string false "Название книги"
// @Param author query string false "Автор"
// @Param genre query []string false "Жанры (можно несколько)" collectionFormat(multi)
// @Success 200 {array} domain.Book
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/search [get]
func (h *BookHandler) SearchBooks(c *gin.Context) {
	filter := domain.BookFilter{
		Title:  c.Query("title"),
		Author: c.Query("author"),
		Genres: c.QueryArray("genre"),
	}

	books, err := h.bookUC.SearchBooks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// CountBooks godoc
// @Summary Подсчитать общее количество книг
// @Tags books
// @Produce json
// @Success 200 {object} map[string]int64
// @Failure 500 {object} map[string]string
// @Router /books/count [get]
func (h *BookHandler) CountBooks(c *gin.Context) {
	count, err := h.bookUC.CountBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]int64{"count": count})
}
