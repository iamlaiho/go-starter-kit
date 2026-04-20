package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv              string   `mapstructure:"APP_ENV"`
	Port                string   `mapstructure:"PORT"`
	JWTSecret           string   `mapstructure:"JWT_SECRET"`
	MaxRequestBodyBytes int64    `mapstructure:"MAX_REQUEST_BODY_BYTES"`
	CORSAllowedOrigins  []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MAX_REQUEST_BODY_BYTES", 1048576)

	if err := viper.BindEnv("CORS_ALLOWED_ORIGINS"); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET is required")
	}

	if cfg.AppEnv == "production" && len(cfg.CORSAllowedOrigins) == 0 {
		return nil, fmt.Errorf("config: CORS_ALLOWED_ORIGINS is required in production")
	}

	if cfg.AppEnv != "production" && len(cfg.CORSAllowedOrigins) == 0 {
		cfg.CORSAllowedOrigins = []string{"*"}
	}

	// CORS_ALLOWED_ORIGINS may arrive as a comma-separated string from env
	if len(cfg.CORSAllowedOrigins) == 1 && strings.Contains(cfg.CORSAllowedOrigins[0], ",") {
		cfg.CORSAllowedOrigins = strings.Split(cfg.CORSAllowedOrigins[0], ",")
	}

	return &cfg, nil
}
