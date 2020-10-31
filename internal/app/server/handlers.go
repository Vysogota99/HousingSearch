package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"context"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"github.com/gin-gonic/gin"
)

const (
	INVALID_REQUEST_BODY  = "Веправильное тело запроса"
	INTERNAL_SERVER_ERROR = "Внутренняя ошибка сервера"
	ALREADY_EXISTS        = "Невозможно создать запись, она уже существует"
	LOT_NOT_FOUND         = "Квартира или комната не найдена"
	BAD_ORDERBY_PARAMS    = "Параметры для сортировки установлены неправильно"
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
		Name            string `json:"name" binding:"required"`
		LastName        string `json:"lastname" binding:"required"`
		Sex             string `json:"sex" binding:"required"`
		DateOfBirth     string `json:"dateOfBirth" binding:"required"`
		Password        string `json:"password" binding:"required"`
		TelephoneNumber string `json:"telephoneNumber" binding:"required"`
		Role            string `json:"role" binding:"required"`
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

// PostLotHandler - добавляет новый лот(квартиру) в систему
func (r *Router) PostLotHandler(c *gin.Context) {
	lot := &models.Lot{}
	if err := c.ShouldBindJSON(lot); err != nil {
		respond(c, http.StatusUnprocessableEntity, INVALID_REQUEST_BODY, err.Error())
		return
	}

	uID := 1
	lot.OwnerID = uID
	if err := r.store.Lot().Create(context.Background(), lot); err != nil {
		respond(c, http.StatusOK, ALREADY_EXISTS, err.Error())
		return
	}
	respond(c, http.StatusCreated, lot, "")
}

// GetLotsHandler - получить список квартир
func (r *Router) GetLotsHandler(c *gin.Context) {
	offset := c.DefaultQuery("offset", "1")
	limit := c.DefaultQuery("limit", "10")
	orderStr := c.DefaultQuery("order_by", "created_at desc")

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	orderBy := strings.Split(orderStr, " ")
	if cap(orderBy) != 2 {
		respond(c, http.StatusBadRequest, nil, BAD_ORDERBY_PARAMS)
	}

	res, err := r.store.Lot().GetFlats(context.Background(), limitInt, offsetInt, nil, orderBy)
	if err != nil {
		respond(c, http.StatusOK, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, res, "")
}

// GetLotHandler - получить полную информацию о конкретном лоте
func (r *Router) GetLotHandler(c *gin.Context) {
	lotID := c.Param("lotid")
	lotIDInt, err := strconv.Atoi(lotID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	res, err := r.store.Lot().GetFlat(context.Background(), lotIDInt)
	switch {
	case err == sql.ErrNoRows:
		respond(c, http.StatusNotFound, res, LOT_NOT_FOUND)
		return
	case err != nil:
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, res, "")
}

// GetRoomsMapHandler - получить список комнат, чтобы потом разместить их на карте
func (r *Router) GetRoomsMapHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "1000")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	rooms, err := r.store.Room().GetRoomsForMap(context.Background(), limitInt)
	switch {
	case err == sql.ErrNoRows:
		respond(c, http.StatusNotFound, rooms, LOT_NOT_FOUND)
		return
	case err != nil:
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, rooms, "")
}

func respond(c *gin.Context, code int, result interface{}, err string) {
	c.JSON(
		code,
		gin.H{
			"result": result,
			"error":  err,
		},
	)
}
