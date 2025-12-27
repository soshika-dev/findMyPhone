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

// DeviceHandler handles device endpoints.
type DeviceHandler struct {
	uc     *usecase.DeviceUseCase
	logger *zap.Logger
}

// NewDeviceHandler creates a DeviceHandler.
func NewDeviceHandler(uc *usecase.DeviceUseCase, logger *zap.Logger) *DeviceHandler {
	return &DeviceHandler{uc: uc, logger: logger}
}

// RegisterRoutes sets up device routes.
func (h *DeviceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/devices", h.createDevice)
}

func (h *DeviceHandler) createDevice(c *gin.Context) {
	var req dtos.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid device payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device := &domain.Device{
		DeviceID:   req.DeviceID,
		IMEI:       req.IMEI,
		Generation: req.Generation,
		Name:       req.Name,
		Lost:       req.Lost,
	}

	if err := h.uc.CreateDevice(c.Request.Context(), device); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.DeviceResponse{
		ID:         device.ID,
		DeviceID:   device.DeviceID,
		IMEI:       device.IMEI,
		Generation: device.Generation,
		Name:       device.Name,
		Lost:       device.Lost,
	})
}

func (h *DeviceHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "resource conflict"})
	default:
		h.logger.Error("internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
