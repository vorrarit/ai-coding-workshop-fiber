package config

import (
	"os"
)

type Config struct {
	DatabasePath string
	JWTSecret    []byte
	ServerPort   string
	AppName      string
}

func LoadConfig() *Config {
	jwtSecret := []byte("your-secret-key-change-this-in-production")
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}

	databasePath := "users.db"
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		databasePath = dbPath
	}

	serverPort := ":3000"
	if port := os.Getenv("PORT"); port != "" {
		serverPort = ":" + port
	}

	return &Config{
		DatabasePath: databasePath,
		JWTSecret:    jwtSecret,
		ServerPort:   serverPort,
		AppName:      "Fiber API Server v1.0.0",
	}
}
