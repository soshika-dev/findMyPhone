package handlers

import (
	"net/http"
	"time"

	"findMyPhone/internal/domain"
	"findMyPhone/internal/interface/http/dtos"
	"findMyPhone/internal/usecase"

	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogHandler handles log endpoints.
type LogHandler struct {
	uc     *usecase.LogUseCase
	logger *zap.Logger
}

// NewLogHandler constructs a LogHandler.
func NewLogHandler(uc *usecase.LogUseCase, logger *zap.Logger) *LogHandler {
	return &LogHandler{uc: uc, logger: logger}
}

// RegisterRoutes registers log routes.
func (h *LogHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/logs", h.createLog)
	rg.GET("/logs/last/:device_id", h.getLastLog)
}

func (h *LogHandler) createLog(c *gin.Context) {
	var req dtos.CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid log payload", zap.Error(err))
		sendResponse(c, http.StatusBadRequest, "invalid log payload: "+err.Error(), nil)
		return
	}

	logEntry := &domain.Log{
		DeviceID:  req.DeviceID,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	}

	if err := h.uc.CreateLog(c.Request.Context(), logEntry); err != nil {
		h.handleError(c, err)
		return
	}

	response := dtos.LogResponse{
		ID:        logEntry.ID,
		DeviceID:  logEntry.DeviceID,
		Longitude: logEntry.Longitude,
		Latitude:  logEntry.Latitude,
		CreatedAt: logEntry.CreatedAt.Format(time.RFC3339),
	}

	sendResponse(c, http.StatusCreated, "log created successfully", response)
}

func (h *LogHandler) getLastLog(c *gin.Context) {
	deviceID := c.Param("device_id")
	logEntry, err := h.uc.GetLastLogByDeviceID(c.Request.Context(), deviceID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := dtos.LogResponse{
		ID:        logEntry.ID,
		DeviceID:  logEntry.DeviceID,
		Longitude: logEntry.Longitude,
		Latitude:  logEntry.Latitude,
		CreatedAt: logEntry.CreatedAt.Format(time.RFC3339),
	}

	sendResponse(c, http.StatusOK, "last log fetched successfully", response)
}

func (h *LogHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		sendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, domain.ErrNotFound):
		sendResponse(c, http.StatusNotFound, "resource not found", nil)
	default:
		h.logger.Error("internal error", zap.Error(err))
		sendResponse(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
