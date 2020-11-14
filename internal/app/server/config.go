package server

import (
	"fmt"
	"os"
)

// Config ...
type Config struct {
	ServerPort      string
	FrontPort       string
	BaseLocURL      string
	BaseProdURL     string
	AuthServicePort string
	DBConnString    string
	StorageLevel    string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return nil, fmt.Errorf("No SERVER_PORT in .env")
	}

	frontPort, exists := os.LookupEnv("FRONT_PORT")
	if !exists {
		return nil, fmt.Errorf("No FRONT_PORT in .env")
	}

	baseLocURL, exists := os.LookupEnv("BASE_PROD_URL")
	if !exists {
		return nil, fmt.Errorf("No BASE_PROD_URL in .env")
	}

	baseProdURL, exists := os.LookupEnv("BASE_PROD_URL")
	if !exists {
		return nil, fmt.Errorf("No BASE_PROD_URL in .env")
	}

	authServicePort, exists := os.LookupEnv("AUTH_SERVICE_PORT")
	if !exists {
		panic("No AUTH_SERVICE_PORT in .env")
	}

	dbConnString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		return nil, fmt.Errorf("No DB_CONN_STRING in .env")
	}

	storageLevel, exists := os.LookupEnv("STORAGE_LEVEL")
	if !exists {
		return nil, fmt.Errorf("No STORAGE_LEVEL in .env")
	}

	return &Config{
		ServerPort:      serverPort,
		DBConnString:    dbConnString,
		FrontPort:       frontPort,
		BaseLocURL:      baseLocURL,
		BaseProdURL:     baseProdURL,
		AuthServicePort: authServicePort,
		StorageLevel:    storageLevel,
	}, nil
}
