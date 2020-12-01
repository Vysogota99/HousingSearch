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
		api.OPTIONS("/logout", r.HeadersMiddleware(), r.OptionsHandler)
		api.POST("/login", r.HeadersMiddleware(), r.LogInHandler)

		api.POST("/lot", r.HeadersMiddleware(), r.PostLotHandler)
		api.GET("/lot", r.HeadersMiddleware(), r.GetLotsHandler)
		api.GET("/lot/id/:lotid", r.HeadersMiddleware(), r.GetLotHandler)
		api.GET("/lot/owner/ads", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetLotsOwnerHandler)

		api.GET("/rooms", r.HeadersMiddleware(), r.GetRoomsHandler)
		api.GET("/rooms/id/:roomid", r.HeadersMiddleware(), r.GetRoomHandler)
		api.GET("/rooms/owner/ads", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetRoomsOwnerHandler)
		api.POST("/rooms", r.HeadersMiddleware(), r.PostRoomHandler)
		api.DELETE("/rooms/room/:roomid", r.HeadersMiddleware(), r.DeleteRoomHandler)
		api.PATCH("/rooms/room/:roomid", r.HeadersMiddleware(), r.UpdateRoomHandler)
		api.DELETE("/rooms/living_place/:lpid", r.HeadersMiddleware(), r.DeleteLivingPlaceHandler)
		api.PATCH("/rooms/living_place/:lpid", r.HeadersMiddleware(), r.UpdateLivingPlaceHandler)
	}

	return r.router
}
