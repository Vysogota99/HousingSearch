package ginrouter

import (
	"fmt"
	"net/http"

	"context"

	"github.com/Vysogota99/school/pkg/authService"
	"github.com/gin-gonic/gin"
)

// TestAPIHandler - handle request from outside to check accessibility of the server
func (r *GinRouter) TestAPIHandler(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

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

// LoginHandler - ...
func (r *GinRouter) LoginHandler(c *gin.Context) {
	res, err := r.authClient.CreateAuth(context.Background(), &authService.CreateAuthRequest{ID: "Ivan"})
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

	c.JSON(http.StatusOK, gin.H{"result": res})
	return
}
