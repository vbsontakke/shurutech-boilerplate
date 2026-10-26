// Package config loads application configuration from environment variables.
// A local .env file is loaded if present; real environment variables win.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	Port            string
	LogLevel        string
	CORSOrigins     []string
	ShutdownTimeout time.Duration
	DB              DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int32
}

// DSN builds a postgres connection URL with credentials safely escaped.
func (d DBConfig) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     d.Host + ":" + d.Port,
		Path:     "/" + d.Name,
		RawQuery: "sslmode=" + url.QueryEscape(d.SSLMode),
	}
	return u.String()
}

// Load reads config from the environment (and .env if present).
func Load() (Config, error) {
	_ = godotenv.Load() // optional; ignore "file not found"

	maxConns, err := strconv.Atoi(getEnv("DB_MAX_CONNS", "10"))
	if err != nil {
		return Config{}, fmt.Errorf("DB_MAX_CONNS must be an integer: %w", err)
	}

	cfg := Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("APP_PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		CORSOrigins:     splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173")),
		ShutdownTimeout: 10 * time.Second,
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: int32(maxConns),
		},
	}

	var missing []string
	if cfg.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
