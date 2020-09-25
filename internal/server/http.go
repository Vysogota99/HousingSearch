package server

import (
	"github.com/Vysogota99/school/internal/server/routers/ginrouter"
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
	ginRouter := ginrouter.NewRouter(conf.ServerPort)
	return &HTTP{
		router: ginRouter,
		conf:   conf,
	}
}
