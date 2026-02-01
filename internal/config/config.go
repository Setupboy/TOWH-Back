package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Jwt      JwtConfig
	Upload   UploadConfig
	SMTP     SMTPConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JwtConfig struct {
	Secret              string
	ExpireIn            time.Duration
	RefreshTokenExpires time.Duration
}
type UploadConfig struct {
	Path           string
	MaxFileSize    int64
	UploadProvider string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpires, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "720h"))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8000"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Jwt: JwtConfig{
			Secret:              getEnv("JWT_SECRET", "your_jwt_secret_key"),
			ExpireIn:            jwtExpiresIn,
			RefreshTokenExpires: refreshTokenExpires,
		},
		Upload: UploadConfig{
			Path:           getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:    maxUploadSize,
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@shop.com"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
