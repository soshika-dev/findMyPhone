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
	rg.PUT("/devices/:device_id", h.updateDevice)
}

func (h *DeviceHandler) createDevice(c *gin.Context) {
	var req dtos.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid device payload", zap.Error(err))
		sendResponse(c, http.StatusBadRequest, "invalid device payload: "+err.Error(), nil)
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

	response := dtos.DeviceResponse{
		ID:         device.ID,
		DeviceID:   device.DeviceID,
		IMEI:       device.IMEI,
		Generation: device.Generation,
		Name:       device.Name,
		Lost:       device.Lost,
	}

	sendResponse(c, http.StatusCreated, "device created successfully", response)
}

func (h *DeviceHandler) updateDevice(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req dtos.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid device payload", zap.Error(err))
		sendResponse(c, http.StatusBadRequest, "invalid device payload: "+err.Error(), nil)
		return
	}

	device := &domain.Device{
		IMEI:       req.IMEI,
		Generation: req.Generation,
		Name:       req.Name,
		Lost:       req.Lost,
	}

	updated, err := h.uc.UpdateDevice(c.Request.Context(), deviceID, device)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dtos.DeviceResponse{
		ID:         updated.ID,
		DeviceID:   updated.DeviceID,
		IMEI:       updated.IMEI,
		Generation: updated.Generation,
		Name:       updated.Name,
		Lost:       updated.Lost,
	}

	sendResponse(c, http.StatusOK, "device updated successfully", response)
}

func (h *DeviceHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		sendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, domain.ErrConflict):
		sendResponse(c, http.StatusConflict, "resource conflict", nil)
	default:
		h.logger.Error("internal error", zap.Error(err))
		sendResponse(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
