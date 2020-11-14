package server

import (
	"strconv"

	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/Vysogota99/HousingSearch/internal/app/store/postgres"
	"github.com/Vysogota99/HousingSearch/pkg/authService"
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
	storageLevel, _ := strconv.Atoi(conf.StorageLevel)
	return &Server{
		Conf:  conf,
		Store: postgres.New(conf.DBConnString, storageLevel),
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
	s.router.Setup().Run(s.Conf.ServerPort)
	return nil
}

func (s *Server) initRouter() {
	router := NewRouter(s.Conf.ServerPort, s.Store, s.authClient)
	s.router = router
}
