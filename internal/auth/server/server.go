package server

import (
	"log"
	"net"

	"github.com/Vysogota99/school/internal/auth/store"
	"github.com/Vysogota99/school/internal/auth/store/postgresstore"
	"github.com/Vysogota99/school/pkg/authService"
	"google.golang.org/grpc"
)

// GRPCServer ...
type GRPCServer struct {
	Conf  *Config
	Store store.Store
}

// NewGRPCServer ...
func NewGRPCServer(conf *Config) (*GRPCServer, error) {
	server := &GRPCServer{}
	server.Store = postgresstore.New(conf.DbConString)
	server.Conf = conf
	return server, nil
}

// Start - start tcp server
func (s *GRPCServer) Start() error {
	srv := grpc.NewServer()
	authService.RegisterAuthorizerServer(srv, s)
	listenerTCP, err := net.Listen("tcp", s.Conf.AuthServicePort)
	if err != nil {
		return err
	}

	log.Printf("Starting tcp at %s\n", s.Conf.AuthServicePort)
	if err := srv.Serve(listenerTCP); err != nil {
		return err
	}

	return nil
}
