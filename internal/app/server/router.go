package server

import (
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"github.com/gin-gonic/gin"
)

var (
	sessionName = "auth"
	adminRole   = 1
	userRile    = 2
)

// Router ...
type Router struct {
	router     *gin.Engine
	serverPort string
	store      store.Store
	authClient authService.AuthorizerClient
}

// NewRouter - helper for initialization http
func NewRouter(serverPort string, store store.Store, authClient authService.AuthorizerClient) *Router {
	return &Router{
		router:     gin.Default(),
		serverPort: serverPort,
		store:      store,
		authClient: authClient,
	}
}

// Setup - найстройка роутера
func (r *Router) Setup() *gin.Engine {
	api := r.router.Group("/api")
	{
		api.GET("/", r.HeadersMiddleware(), r.TestAPIHandler)
		api.POST("/", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.TestAPIHandler)

		api.POST("/signup", r.HeadersMiddleware(), r.SignUPHandler)
		api.POST("/logout", r.HeadersMiddleware(), r.LogoutHandler)
		api.POST("/login", r.HeadersMiddleware(), r.LogInHandler)
	}
	r.router.Run(r.serverPort)

	return r.router
}
