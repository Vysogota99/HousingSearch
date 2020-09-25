package server

// Server ...
type Server struct {
	Conf *Config
	HTTP *HTTP
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	return &Server{
		Conf: conf,
		HTTP: NewHTTP(conf),
	}, nil
}

// Start - start the server
func (s *Server) Start() {
	s.HTTP.router.Run()
}
