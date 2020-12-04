package server

import (
	"os"
)

// Config ...
type Config struct {
	AuthServicePort   string
	DbConString       string
	RedisDSN          string
	JWTAccessExpTime  string
	JWTRefreshExpTime string
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

	jwtAccessExpTime, exists := os.LookupEnv("JWT_ACCESS_EXPIRE_TIME")
	if !exists {
		panic("No JWT_ACCESS_EXPIRE_TIME in .env")
	}

	jwtRefreshExpTime, exists := os.LookupEnv("JWT_ACCESS_REFRESH_TIME")
	if !exists {
		panic("No JWT_ACCESS_REFRESH_TIME in .env")
	}

	return &Config{
		AuthServicePort:   authServicePort,
		DbConString:       dbConString,
		RedisDSN:          redisDSN,
		JWTAccessExpTime:  jwtAccessExpTime,
		JWTRefreshExpTime: jwtRefreshExpTime,
	}, nil
}
