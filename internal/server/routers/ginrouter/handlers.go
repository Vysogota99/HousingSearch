package ginrouter

import (
	"net/http"

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
