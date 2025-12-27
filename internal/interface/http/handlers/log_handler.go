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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusCreated, dtos.LogResponse{
		ID:        logEntry.ID,
		DeviceID:  logEntry.DeviceID,
		Longitude: logEntry.Longitude,
		Latitude:  logEntry.Latitude,
		CreatedAt: logEntry.CreatedAt.Format(time.RFC3339),
	})
}

func (h *LogHandler) getLastLog(c *gin.Context) {
	deviceID := c.Param("device_id")
	logEntry, err := h.uc.GetLastLogByDeviceID(c.Request.Context(), deviceID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.LogResponse{
		ID:        logEntry.ID,
		DeviceID:  logEntry.DeviceID,
		Longitude: logEntry.Longitude,
		Latitude:  logEntry.Latitude,
		CreatedAt: logEntry.CreatedAt.Format(time.RFC3339),
	})
}

func (h *LogHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	default:
		h.logger.Error("internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
