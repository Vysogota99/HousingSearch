package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware - проверяет валидность токена
func (r *Router) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := c.Request.Header.Get("Authorization")
		log.Println(jwt)
		log.Println(c.Request.Header)
		res, err := r.authClient.CheckAuthUser(context.Background(), &authService.CheckAuthRequest{Jwt: jwt})
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"message": fmt.Errorf("Ошибка при запросе к сервису авторизации").Error(),
					"error":   err.Error(),
				},
			)
			c.Abort()
			return
		}

		c.Set("telephone_number", res.TelephoneNumber)
		c.Set("user_id", res.UserID)
		c.Next()
	}
}

// HeadersMiddleware - устанавливает заголовки
func (r *Router) HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Add("Content-Type", "application/json")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Add("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Add("Access-Control-Request-Headers", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "authorization")

		c.Next()
	}
}
