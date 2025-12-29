package handlers

import (
	"github.com/gin-gonic/gin"

	"findMyPhone/internal/interface/http/dtos"
)

func sendResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, dtos.Response{ //nolint: exhaustive struct - clarity prioritized
		Data:    data,
		Status:  status,
		Message: message,
	})
}
