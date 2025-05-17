package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	customErr "library-Mongo/internal/errors"
	"library-Mongo/internal/usecase"
	"library-Mongo/internal/usecase/dto"
	"net/http"
	"time"
)

type BorrowHandler struct {
	borrowUC usecase.BorrowUC
}

func NewBorrowHandler(borrowUC usecase.BorrowUC) *BorrowHandler {
	return &BorrowHandler{borrowUC: borrowUC}
}

// BorrowBook godoc
// @Summary Выдача книги
// @Tags borrow
// @Accept json
// @Produce json
// @Param input body dto.BorrowBookInput true "Данные для выдачи"
// @Success 200 {object} domain.Borrow
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow [post]
func (h *BorrowHandler) BorrowBook(c *gin.Context) {
	var input dto.BorrowBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	borrow, err := h.borrowUC.BorrowBook(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		case errors.Is(err, customErr.ErrBookAlreadyBorrowed):
			c.JSON(http.StatusBadRequest, gin.H{"error": "book already borrowed"})
		case errors.Is(err, customErr.ErrBookNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		case errors.Is(err, customErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, borrow)
}

// ReturnBook godoc
// @Summary Возврат книги
// @Tags borrow
// @Accept json
// @Produce json
// @Param input body dto.ReturnBookInput true "Данные для возврата"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow/return [post]
func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	var input dto.ReturnBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.borrowUC.ReturnBook(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		case errors.Is(err, customErr.ErrBorrowNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "borrow not found"})
		case errors.Is(err, customErr.ErrAlreadyReturned):
			c.JSON(http.StatusBadRequest, gin.H{"error": "book already returned"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "returned"})
}

// GetBorrowHistory godoc
// @Summary История выдач пользователя
// @Tags borrow
// @Produce json
// @Param userID path string true "ID пользователя"
// @Success 200 {object} dto.BorrowHistoryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow/history/{userID} [get]
func (h *BorrowHandler) GetBorrowHistory(c *gin.Context) {
	userID := c.Param("userID")
	history, err := h.borrowUC.GetBorrowHistory(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, history)
}

// GetOverdueBorrows godoc
// @Summary Просроченные книги
// @Tags borrow
// @Produce json
// @Success 200 {array} dto.OverdueReportItem
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow/overdue [get]
func (h *BorrowHandler) GetOverdueBorrows(c *gin.Context) {
	result, err := h.borrowUC.GetOverdueBorrows(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetDailyBorrowStats godoc
// @Summary График нагрузки (уникальные читатели)
// @Tags borrow
// @Produce json
// @Param from query string true "Дата начала (YYYY-MM-DD)"
// @Param to query string true "Дата конца (YYYY-MM-DD)"
// @Success 200 {array} domain.BorrowStat
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow/stats [get]
func (h *BorrowHandler) GetDailyBorrowStats(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
		return
	}

	stats, err := h.borrowUC.GetDailyBorrowStats(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// CountActiveBorrows godoc
// @Summary Кол-во активных выдач
// @Tags borrow
// @Produce json
// @Success 200 {object} dto.CountResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /borrow/active-count [get]
func (h *BorrowHandler) CountActiveBorrows(c *gin.Context) {
	count, err := h.borrowUC.CountActiveBorrows(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, map[string]int64{"count": count})
}
