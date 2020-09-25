package auth

import (
	"context"

	"github.com/Vysogota99/school/pkg/authService"
)

// GRPCServer ...
type GRPCServer struct {
}

// CreateAuth ...
func (s *GRPCServer) CreateAuth(ctx context.Context, req *authService.CreateAuthRequest) (*authService.CreateAutResponse, error) {
	return &authService.CreateAutResponse{
		TDetails: nil,
		Error:    req.ID,
	}, nil
}
