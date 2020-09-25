package server

import (
	"github.com/Vysogota99/school/internal/app/server/routers/ginrouter"
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
func NewHTTP(conf *Config) *HTTP {
	return &HTTP{
		router: ginrouter.NewRouter(conf.ServerPort),
		conf:   conf,
	}
}
