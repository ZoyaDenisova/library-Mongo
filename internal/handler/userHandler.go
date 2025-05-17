package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"library-Mongo/internal/domain"
	customErr "library-Mongo/internal/errors"
	"library-Mongo/internal/usecase"
	"library-Mongo/internal/usecase/dto"
	"net/http"
)

type UserHandler struct {
	userUC usecase.UserUC
}

func NewUserHandler(userUC usecase.UserUC) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body dto.RegisterUserInput true "Данные пользователя"
// @Success 200 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input dto.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}
	user, err := h.userUC.RegisterUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Login godoc
// @Summary Аутентификация пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Телефон и пароль"
// @Success 200 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}
	user, err := h.userUC.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		case errors.Is(err, customErr.ErrUserBlocked):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "user is blocked"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userUC.GetUserByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid ID"})
		case errors.Is(err, customErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// SearchUsers godoc
// @Summary Поиск пользователей
// @Tags users
// @Produce json
// @Param fullName query string false "ФИО"
// @Param phone query string false "Телефон"
// @Param role query string false "Роль"
// @Param onlyActive query boolean false "Только активные"
// @Success 200 {array} domain.User
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	q := c.Query("query")

	filter := domain.UserFilter{
		FullNameContains: q,
		Phone:            q,
		Role:             q,
	}
	if activeStr := c.Query("onlyActive"); activeStr != "" {
		val := activeStr == "true"
		filter.OnlyActive = &val
	}
	users, err := h.userUC.SearchUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Обновление пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body dto.UpdateUserInput true "Данные обновления"
// @Success 200 {object} dto.StatusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}
	if err := h.userUC.UpdateUser(c.Request.Context(), input); err != nil {
		switch {
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid ID"})
		case errors.Is(err, customErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "updated"})
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} dto.StatusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUC.DeleteUser(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, customErr.ErrInvalidID):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid ID"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "deleted"})
}

// BlockUser / UnblockUser — можешь оформить аналогично по схеме UpdateUser.
