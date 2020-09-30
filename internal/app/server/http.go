package server

import (
	"fmt"

	"github.com/Vysogota99/school/internal/app/server/routers/ginrouter"
	"github.com/Vysogota99/school/pkg/authService"
	"google.golang.org/grpc"
)

// HTTP ...
type HTTP struct {
	router Router
	conf   *Config
}

// Router - interface for router
type Router interface {
	Run()
}

// NewHTTP - helper for initialization http
func NewHTTP(conf *Config) (*HTTP, error) {
	authServiceClient, err := getAuthSrerviceClient(conf.AuthServicePort)
	if err != nil {
		return nil, fmt.Errorf("Can't dial to authService. %w", err)
	}

	return &HTTP{
		router: ginrouter.NewRouter(conf.ServerPort, authServiceClient),
		conf:   conf,
	}, nil
}

func getAuthSrerviceClient(port string) (authService.AuthorizerClient, error) {
	connection, err := grpc.Dial(port, grpc.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("Can't dial to authService. %w", err)
	}

	return authService.NewAuthorizerClient(connection), nil
}
