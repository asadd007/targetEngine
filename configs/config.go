package configs

import (
	"os"
	"strconv"
)

type Config struct {
	Port              int
	LogLevel          string
	EnableMetrics     bool
	MetricsPort       int
	EnableHealthCheck bool
	Database          struct {
		PostgresURI string
	}
}

func NewConfig() *Config {
	cfg := &Config{
		Port:              8080,
		LogLevel:          "info",
		EnableMetrics:     false,
		MetricsPort:       9090,
		EnableHealthCheck: true,
	}
	cfg.Database.PostgresURI = "postgres://postgres:postgres@localhost:5432/targeting_engine?sslmode=disable"
	return cfg
}

func (c *Config) LoadFromEnv() {
	// Server settings
	if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port > 0 {
		c.Port = port
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}

	if enableMetrics := os.Getenv("ENABLE_METRICS"); enableMetrics != "" {
		c.EnableMetrics, _ = strconv.ParseBool(enableMetrics)
	}

	if metricsPort, err := strconv.Atoi(os.Getenv("METRICS_PORT")); err == nil && metricsPort > 0 {
		c.MetricsPort = metricsPort
	}

	if enableHealthCheck := os.Getenv("ENABLE_HEALTH_CHECK"); enableHealthCheck != "" {
		c.EnableHealthCheck, _ = strconv.ParseBool(enableHealthCheck)
	}

	// Database settings
	if postgresURI := os.Getenv("POSTGRES_URI"); postgresURI != "" {
		c.Database.PostgresURI = postgresURI
	}
}
