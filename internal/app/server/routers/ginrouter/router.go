package ginrouter

import (
	"github.com/Vysogota99/school/pkg/authService"
	"github.com/gin-gonic/gin"
)

// GinRouter - router base on gin
type GinRouter struct {
	router     *gin.Engine
	authClient authService.AdderClient
	port       string
}

// NewRouter - helper fo initialozation ginRouter
func NewRouter(port string, authClient authService.AdderClient) *GinRouter {
	return &GinRouter{
		port:       port,
		router:     gin.Default(),
		authClient: authClient,
	}
}

// Run - run the router
func (r *GinRouter) Run() {
	api := r.router.Group("/api")
	{
		api.GET("/", r.TestAPIHandler)
		api.POST("/", r.TestAPIHandler)

		api.POST("/login", r.LoginHandler)
	}
	r.router.Run(r.port)
}
