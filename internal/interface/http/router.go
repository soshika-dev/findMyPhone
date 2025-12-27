package http

import (
	"findMyPhone/internal/interface/http/handlers"
	"findMyPhone/internal/interface/http/middleware"
	"findMyPhone/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter sets up routes and middleware.
func NewRouter(userUC *usecase.UserUseCase, deviceUC *usecase.DeviceUseCase, logUC *usecase.LogUseCase, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	middleware.Setup(r)

	v1 := r.Group("/api/v1")

	userHandler := handlers.NewUserHandler(userUC, logger)
	userHandler.RegisterRoutes(v1)

	deviceHandler := handlers.NewDeviceHandler(deviceUC, logger)
	deviceHandler.RegisterRoutes(v1)

	logHandler := handlers.NewLogHandler(logUC, logger)
	logHandler.RegisterRoutes(v1)

	return r
}
