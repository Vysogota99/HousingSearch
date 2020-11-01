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
	INVALID_REQUEST_BODY  = "Неправильное тело запроса"
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

// GetRoomsHandler - получить список комнат, чтобы потом разместить их на карте
func (r *Router) GetRoomsHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "1000")
	offset := c.DefaultQuery("offset", "1")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	fieldsToFilter := mapOfParams(c, []string{"maxresidents", "currnumberofresidents", "numofwindows", "balcony", "numoftables", "numofchairs", "tv", "numofcupboards", "area"})
	rooms, err := r.store.Room().GetRooms(context.Background(), limitInt, offsetInt, fieldsToFilter)
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

// GetRoomHandler - возвращает информацию о комнате
func (r *Router) GetRoomHandler(c *gin.Context) {
	roomID := c.Param("roomid")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	room, err := r.store.Room().GetRoom(context.Background(), roomIDInt)
	switch {
	case err == sql.ErrNoRows:
		respond(c, http.StatusNotFound, room, LOT_NOT_FOUND)
		return
	case err != nil:
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	lpArray, err := r.store.Room().GetLivingPlace(context.Background(), roomIDInt)
	switch {
	case err == sql.ErrNoRows:
		respond(c, http.StatusNotFound, lpArray, LOT_NOT_FOUND)
		return
	case err != nil:
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	room.LivingPlaces = lpArray
	respond(c, http.StatusOK, room, "")
}

// PostRoomHandler ...
func (r *Router) PostRoomHandler(c *gin.Context) {
	room := models.Room{}
	if err := c.ShouldBindJSON(&room); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	if room.FlatID == 0 {
		respond(c, http.StatusUnprocessableEntity, nil, INVALID_REQUEST_BODY)
		return
	}

	if err := r.store.Room().Create(context.Background(), &room); err != nil {
		respond(c, http.StatusInternalServerError, nil, INTERNAL_SERVER_ERROR)
		return
	}

	respond(c, http.StatusCreated, room, "")
}

// DeleteRoomHandler - удаляет комнату
func (r *Router) DeleteRoomHandler(c *gin.Context) {
	roomID := c.Param("roomid")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := r.store.Room().DeleteRoom(context.Background(), roomIDInt); err != nil {
		respond(c, http.StatusInternalServerError, nil, INTERNAL_SERVER_ERROR)
		return
	}
	respond(c, http.StatusOK, nil, "")
}

// DeleteLivingPlaceHandler - удаляет комнату
func (r *Router) DeleteLivingPlaceHandler(c *gin.Context) {
	lpID := c.Param("lpid")
	lpIDInt, err := strconv.Atoi(lpID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := r.store.Room().DeleteLivingPlace(context.Background(), lpIDInt); err != nil {
		respond(c, http.StatusInternalServerError, nil, INTERNAL_SERVER_ERROR)
		return
	}
	respond(c, http.StatusOK, nil, "")
}

// UpdateRoomHandler - обновляет данные о комнате
func (r *Router) UpdateRoomHandler(c *gin.Context) {
	roomID := c.Param("roomid")
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	type request struct {
		Fields map[string]interface{} `json:"fields" binding:"required"`
	}
	expRequest := &request{}
	if err := c.ShouldBindJSON(&expRequest); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	if err := r.store.Room().UpdateRoom(context.Background(), roomIDInt, expRequest.Fields); err != nil {
		respond(c, http.StatusInternalServerError, nil, INTERNAL_SERVER_ERROR)
		return
	}

	respond(c, http.StatusOK, nil, "")
}

// UpdateLivingPlaceHandler - ...
func (r *Router) UpdateLivingPlaceHandler(c *gin.Context) {
	lpID := c.Param("lpid")
	lpIDInt, err := strconv.Atoi(lpID)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	type request struct {
		Fields map[string]interface{} `json:"fields" binding:"required"`
	}
	expRequest := &request{}
	if err := c.ShouldBindJSON(&expRequest); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	if err := r.store.Room().UpdateLivingPlace(context.Background(), lpIDInt, expRequest.Fields); err != nil {
		respond(c, http.StatusInternalServerError, nil, INTERNAL_SERVER_ERROR)
		return
	}

	respond(c, http.StatusOK, nil, "")
}

//====================================HELPERS================================================================

func respond(c *gin.Context, code int, result interface{}, err string) {
	c.JSON(
		code,
		gin.H{
			"result": result,
			"error":  err,
		},
	)
}

func mapOfParams(c *gin.Context, fields []string) map[string]string {
	result := make(map[string]string)
	for _, field := range fields {
		if value := c.Query(field); value != "" {
			result[field] = value
		}
	}

	return result
}
