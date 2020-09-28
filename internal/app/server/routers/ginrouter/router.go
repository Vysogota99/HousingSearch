package ginrouter

import (
	"github.com/Vysogota99/school/pkg/authService"
	"github.com/gin-gonic/gin"
)

// GinRouter - router base on gin
type GinRouter struct {
	router     *gin.Engine
	authClient authService.AuthorizerClient
	port       string
}

// NewRouter - helper fo initialozation ginRouter
func NewRouter(port string, authClient authService.AuthorizerClient) *GinRouter {
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
		api.POST("/", r.TokenAuthMiddleware(), r.TestAPIHandler)

		api.POST("/signup", r.SignUPHandler)
		api.POST("/logout", r.LogoutHandler)
	}
	r.router.Run(r.port)
}
