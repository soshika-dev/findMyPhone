package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"findMyPhone/internal/infrastructure/config"
	"findMyPhone/internal/infrastructure/db"
	"findMyPhone/internal/infrastructure/logger"
	infraRepo "findMyPhone/internal/infrastructure/repository"
	httpInterface "findMyPhone/internal/interface/http"
	"findMyPhone/internal/usecase"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	log, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}
	defer log.Sync()

	database, err := db.NewGorm(cfg)
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	userRepo := infraRepo.NewUserRepository(database)
	deviceRepo := infraRepo.NewDeviceRepository(database)
	logRepo := infraRepo.NewLogRepository(database)

	userUC := usecase.NewUserUseCase(userRepo)
	deviceUC := usecase.NewDeviceUseCase(deviceRepo)
	logUC := usecase.NewLogUseCase(logRepo, deviceRepo)

	router := httpInterface.NewRouter(userUC, deviceUC, logUC, log)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		log.Info("server starting", zap.String("addr", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
	log.Info("server exited")
}
