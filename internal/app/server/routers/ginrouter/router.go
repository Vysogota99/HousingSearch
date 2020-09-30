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
		api.GET("/", r.HeadersMiddleware(), r.TestAPIHandler)
		api.POST("/", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.TestAPIHandler)

		api.POST("/signup", r.HeadersMiddleware(), r.SignUPHandler)
		api.POST("/logout", r.HeadersMiddleware(), r.LogoutHandler)
		api.POST("/login", r.HeadersMiddleware(), r.LogInHandler)
	}
	r.router.Run(r.port)
}
