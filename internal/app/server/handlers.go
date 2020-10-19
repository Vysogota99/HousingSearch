package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"context"

	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"github.com/gin-gonic/gin"
)

// TestAPIHandler - handle request from outside to check accessibility of the server
func (r *Router) TestAPIHandler(c *gin.Context) {
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
func (r *Router) SignUPHandler(c *gin.Context) {
	type User struct {
		Name            string `JSON:"name" binding:"required"`
		LastName        string `JSON:"lastname" binding:"required"`
		Sex             string `JSON:"sex" binding:"required"`
		DateOfBirth     string `JSON:"dateOfBirth" binding:"required"`
		Password        string `JSON:"password" binding:"required"`
		TelephoneNumber string `JSON:"telephoneNumber" binding:"required"`
		Role            string `JSON:"role" binding:"required"`
	}

	user := &User{}
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

	userRequest := &authService.User{}
	userRequest.PassName = user.Name
	userRequest.PassLastName = user.LastName
	userRequest.PassSex = user.Sex
	userRequest.PassDateOfBirth = user.DateOfBirth
	userRequest.Password = user.Password
	userRequest.TelephoneNumber = user.TelephoneNumber
	role, err := strconv.ParseInt(user.Role, 16, 32)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": fmt.Errorf("Внутренняя ошибка").Error(),
				"error":   err.Error(),
			},
		)
		return
	}
	userRequest.Role = int32(role)

	res, err := r.authClient.SignupUser(context.Background(), &authService.SignUPUserRequest{User: userRequest})
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
			"user":          res.User,
		},
	)
}

// LogoutHandler - выход пользователя из системы
func (r *Router) LogoutHandler(c *gin.Context) {
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
func (r *Router) LogInHandler(c *gin.Context) {
	type User struct {
		TelephoneNumber string `JSON:"telephoneNumber" binding:"required"`
		Password        string `JSON:"password" binding:"required"`
	}

	user := &User{}
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

	userRequest := &authService.User{}
	userRequest.Password = user.Password
	userRequest.TelephoneNumber = user.TelephoneNumber

	res, err := r.authClient.LoginUser(context.Background(), &authService.LoginUserRequest{User: userRequest})
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
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
			"user":          res.User,
		},
	)
}
