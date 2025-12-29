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
        rg.POST("/users/by-device/:device_id", h.updateUser)
        rg.POST("/users/:device_id", h.updateUser)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var req dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid user payload", zap.Error(err))
		sendResponse(c, http.StatusBadRequest, "invalid user payload: "+err.Error(), nil)
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

	response := dtos.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		DeviceID:    user.DeviceID,
		Phone:       user.Phone,
		BackupPhone: user.BackupPhone,
	}

	sendResponse(c, http.StatusCreated, "user created successfully", response)
}

func (h *UserHandler) getUserByDeviceID(c *gin.Context) {
	deviceID := c.Param("device_id")
	user, err := h.uc.GetUserByDeviceID(c.Request.Context(), deviceID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dtos.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		DeviceID:    user.DeviceID,
		Phone:       user.Phone,
		BackupPhone: user.BackupPhone,
	}

	sendResponse(c, http.StatusOK, "user fetched successfully", response)
}

func (h *UserHandler) updateUser(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid user payload", zap.Error(err))
		sendResponse(c, http.StatusBadRequest, "invalid user payload: "+err.Error(), nil)
		return
	}

	user := &domain.User{
		Name:        req.Name,
		Phone:       req.Phone,
		BackupPhone: req.BackupPhone,
	}

	updated, err := h.uc.UpdateUser(c.Request.Context(), deviceID, user)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dtos.UserResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		DeviceID:    updated.DeviceID,
		Phone:       updated.Phone,
		BackupPhone: updated.BackupPhone,
	}

	sendResponse(c, http.StatusOK, "user updated successfully", response)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		sendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, domain.ErrNotFound):
		sendResponse(c, http.StatusNotFound, "resource not found", nil)
	case errors.Is(err, domain.ErrConflict):
		sendResponse(c, http.StatusConflict, "resource conflict", nil)
	default:
		h.logger.Error("internal error", zap.Error(err))
		sendResponse(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
