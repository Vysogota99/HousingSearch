package server

// Server ...
type Server struct {
	Conf *Config
	HTTP *HTTP
}

// NewServer - helper to init server
func NewServer(conf *Config) (*Server, error) {
	http, err := NewHTTP(conf)
	if err != nil {
		return nil, err
	}

	return &Server{
		Conf: conf,
		HTTP: http,
	}, nil
}

// Start - start the server
func (s *Server) Start() {
	s.HTTP.router.Run()
}
