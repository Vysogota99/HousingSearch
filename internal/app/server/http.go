package server

import (
	"fmt"

	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"google.golang.org/grpc"
)

func getAuthSrerviceClient(port string) (authService.AuthorizerClient, error) {
	connection, err := grpc.Dial(port, grpc.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("Can't dial to authService. %w", err)
	}

	return authService.NewAuthorizerClient(connection), nil
}
