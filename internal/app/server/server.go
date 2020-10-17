package server

import (
	"github.com/Vysogota99/school/internal/app/store"
	"github.com/Vysogota99/school/internal/app/store/postgres"
	"github.com/Vysogota99/school/pkg/authService"
)

// Server ...
type Server struct {
	Conf       *Config
	router     *Router
	Store      store.Store
	authClient authService.AuthorizerClient
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	return &Server{
		Conf:  conf,
		Store: postgres.New(conf.DBConnString),
	}, nil
}

// Start - start the server
func (s *Server) Start() error {
	authClient, err := getAuthSrerviceClient(s.Conf.AuthServicePort)
	if err != nil {
		return err
	}

	s.authClient = authClient
	s.initRouter()
	s.router.Setup().Run()
	return nil
}

func (s *Server) initRouter() {
	router := NewRouter(s.Conf.ServerPort, s.Store, s.authClient)
	s.router = router
}
