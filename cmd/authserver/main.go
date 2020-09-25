package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Vysogota99/school/internal/auth"
	"github.com/Vysogota99/school/pkg/authService"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("Error could not find .env file: %w", err))
	}
}

func main() {
	serverPort, exists := os.LookupEnv("AUTH_SERVICE_PORT")
	if !exists {
		panic("No AUTH_SERVICE_PORT in .env")
	}

	server := grpc.NewServer()
	service := &auth.GRPCServer{}
	authService.RegisterAdderServer(server, service)

	l, err := net.Listen("tcp", serverPort)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting authservice on port %s\n", serverPort)
	if err := server.Serve(l); err != nil {
		panic(err)
	}

}
