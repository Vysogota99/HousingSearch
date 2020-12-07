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

		api.POST("/lot", r.HeadersMiddleware(), r.PostLotHandler)
		api.GET("/lot", r.HeadersMiddleware(), r.GetLotsHandler)

		api.PATCH("/lot/update/:lotid", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.UpdateLotHandler)
		api.PATCH("/lot/ad", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.CreteAdHandler)
		api.GET("/lot/id/:lotid/:isconstructot", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetLotHandler)
		api.GET("/lot/owner/ads", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetLotsOwnerHandler)
		api.GET("/lot/owner/construct", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetConstructOwnerHandler)
		api.DELETE("/lot", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.DeleteLotHandler)

		api.GET("/rooms", r.HeadersMiddleware(), r.GetRoomsHandler)
		api.GET("/rooms/id/:roomid", r.HeadersMiddleware(), r.GetRoomHandler)
		api.GET("/rooms/owner/ads", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.GetRoomsOwnerHandler)
		api.POST("/rooms", r.HeadersMiddleware(), r.PostRoomHandler)
		api.DELETE("/rooms/room/:roomid", r.HeadersMiddleware(), r.DeleteRoomHandler)
		api.PATCH("/rooms/room/:roomid", r.TokenAuthMiddleware(), r.HeadersMiddleware(), r.UpdateRoomHandler)
		api.DELETE("/rooms/living_place/:lpid", r.HeadersMiddleware(), r.DeleteLivingPlaceHandler)
		api.PATCH("/rooms/living_place/:lpid", r.HeadersMiddleware(), r.UpdateLivingPlaceHandler)

		api.OPTIONS("/signup", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/logout", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/login", r.HeadersMiddleware(), r.LogInHandler)

		api.OPTIONS("/lot", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/lot/ad", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/lot/id/:lotid/:isconstructot", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/lot/owner/ads", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/lot/owner/construct", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/lot/update/:lotid", r.HeadersMiddleware(), r.OptionsHandler)

		api.OPTIONS("/rooms", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/rooms/id/:roomid", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/rooms/owner/ads", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/rooms/room/:roomid", r.HeadersMiddleware(), r.OptionsHandler)
		api.OPTIONS("/rooms/living_place/:lpid", r.HeadersMiddleware(), r.OptionsHandler)

	}

	return r.router
}
