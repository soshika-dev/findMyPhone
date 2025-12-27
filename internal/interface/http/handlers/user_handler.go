package handlers

import (
	"net/http"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/interface/http/dtos"
	"findMyPhone/internal/usecase"

	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler handles user endpoints.
type UserHandler struct {
	uc     *usecase.UserUseCase
	logger *zap.Logger
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(uc *usecase.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{uc: uc, logger: logger}
}

// RegisterRoutes registers user routes.
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", h.createUser)
	rg.GET("/users/by-device/:device_id", h.getUserByDeviceID)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid user payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Name:        req.Name,
		DeviceID:    req.DeviceID,
		Phone:       req.Phone,
		BackupPhone: req.BackupPhone,
	}

	if err := h.uc.CreateUser(c.Request.Context(), user); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dtos.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		DeviceID:    user.DeviceID,
		Phone:       user.Phone,
		BackupPhone: user.BackupPhone,
	})
}

func (h *UserHandler) getUserByDeviceID(c *gin.Context) {
	deviceID := c.Param("device_id")
	user, err := h.uc.GetUserByDeviceID(c.Request.Context(), deviceID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		DeviceID:    user.DeviceID,
		Phone:       user.Phone,
		BackupPhone: user.BackupPhone,
	})
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "resource conflict"})
	default:
		h.logger.Error("internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
