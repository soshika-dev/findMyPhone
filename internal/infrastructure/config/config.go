package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration loaded from env.
type Config struct {
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	ShutdownGrace time.Duration `mapstructure:"SHUTDOWN_GRACE"`
	DatabaseURL   string        `mapstructure:"DATABASE_URL"`
	DatabaseType  string        `mapstructure:"DATABASE_TYPE"`
}

// Load initializes viper and reads environment variables.
func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetDefault("SERVER_ADDRESS", ":8081")
	viper.SetDefault("SHUTDOWN_GRACE", 15*time.Second)
	viper.SetDefault("DATABASE_TYPE", "postgres")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/findmyphone?sslmode=disable")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
