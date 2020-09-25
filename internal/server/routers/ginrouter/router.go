package ginrouter

import (
	"github.com/gin-gonic/gin"
)

// GinRouter - router base on gin
type GinRouter struct {
	router *gin.Engine
	port   string
}

// NewRouter - helper fo initialozation ginRouter
func NewRouter(port string) *GinRouter {
	return &GinRouter{
		port:   port,
		router: gin.Default(),
	}
}

// Run - run the router
func (r *GinRouter) Run() {
	api := r.router.Group("/api")
	{
		api.GET("/", r.TestAPIHandler)
		api.POST("/", r.TestAPIHandler)
	}
	r.router.Run(r.port)
}
