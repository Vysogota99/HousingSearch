package server

import (
	"os"
)

// Config ...
type Config struct {
	AuthServicePort string
	DbConString     string
	RedisDSN        string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {

	authServicePort, exists := os.LookupEnv("AUTH_SERVICE_PORT")
	if !exists {
		panic("No AUTH_SERVICE_PORT in .env")
	}

	dbConString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		panic("No DB_CONN_STRING in .env")
	}

	redisDSN, exists := os.LookupEnv("REDIS_DSN")
	if !exists {
		panic("No REDIS_DSN in .env")
	}

	return &Config{
		AuthServicePort: authServicePort,
		DbConString:     dbConString,
		RedisDSN:        redisDSN,
	}, nil
}
