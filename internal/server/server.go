package server

// Server ...
type Server struct {
	Conf *Config
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	return &Server{
		Conf: conf,
	}, nil
}

// Start - start the server
func (s *Server) Start() {

}
