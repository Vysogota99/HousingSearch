package main

import (
	"fmt"

	"github.com/Vysogota99/school/internal/auth/server"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("Error could not find .env file: %w", err))
	}
}

func main() {
	conf, err := server.NewConfig()
	if err != nil {
		panic(err)
	}
	serv, err := server.NewGRPCServer(conf)
	if err != nil {
		panic(err)
	}

	if err := serv.Start(); err != nil {
		panic(err)
	}
}
