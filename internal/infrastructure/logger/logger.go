package logger

import "go.uber.org/zap"

// New creates a zap logger with production config.
func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
