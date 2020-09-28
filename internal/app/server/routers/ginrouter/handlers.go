package ginrouter

import (
	"errors"
	"fmt"
	"net/http"

	"context"

	"github.com/Vysogota99/school/pkg/authService"
	"github.com/gin-gonic/gin"
)

// TestAPIHandler - handle request from outside to check accessibility of the server
func (r *GinRouter) TestAPIHandler(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"Body":      c.Request.Body,
			"URL":       c.Request.URL,
			"Method":    c.Request.Method,
			"UserAgent": c.Request.UserAgent(),
			"IP":        c.Request.RemoteAddr,
		},
	)
}

// SignUPHandler - регистрация пользователя
func (r *GinRouter) SignUPHandler(c *gin.Context) {
	user := &authService.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			gin.H{
				"message": errors.New("Неправильное тело запроса").Error(),
				"error":   err.Error(),
			},
		)
		return
	}

	res, err := r.authClient.SignupUser(context.Background(), &authService.SignUPUserRequest{User: user})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": fmt.Errorf("Ошибка при запросе к сервису авторизации").Error(),
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
			"user":          user,
		},
	)
}

// LogoutHandler - выход пользователя из системы
func (r *GinRouter) LogoutHandler(c *gin.Context) {
	jwt := c.Request.Header.Get("Authorization")
	if jwt == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": errors.New("Пользователь не авторизован").Error(),
			},
		)
		return
	}

	res, err := r.authClient.LogOutUser(context.Background(), &authService.LogOutRequest{Jwt: jwt})
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": fmt.Errorf("Ошибка при запросе к сервису авторизации").Error(),
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"result": res.Message,
		},
	)
}

// LogInHandler - автризация пользователя
func (r *GinRouter) LogInHandler(c *gin.Context) {

}