package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	UNAUTH                = "Пользователь не авторизован"
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
	long := c.DefaultQuery("long", "0")
	lat := c.DefaultQuery("lat", "0")
	radius := c.DefaultQuery("radius", "0")

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

	longFl64, err := strconv.ParseFloat(long, 64)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	latFl64, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	radiusInt, err := strconv.Atoi(radius)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	res, err := r.store.Lot().GetFlats(context.Background(), limitInt, offsetInt, nil, false, orderBy, longFl64, latFl64, radiusInt)
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

	isConstructor := c.Param("isconstructot")
	isConstructorBool, err := strconv.ParseBool(isConstructor)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	res, err := r.store.Lot().GetFlatAd(context.Background(), lotIDInt, isConstructorBool)
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
	roomLong := c.DefaultQuery("long", "0")
	roomLat := c.DefaultQuery("lat", "0")
	roomRadius := c.DefaultQuery("radius", "0")

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
	roomLongFl64, err := strconv.ParseFloat(roomLong, 64)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	roomLatFl64, err := strconv.ParseFloat(roomLat, 64)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	roomRadiusInt, err := strconv.Atoi(roomRadius)
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	fieldsToFilter := mapOfParams(c, models.MapRoom)
	rooms, err := r.store.Room().GetRooms(context.Background(), limitInt, offsetInt, fieldsToFilter, false, roomLongFl64, roomLatFl64, roomRadiusInt)
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

// GetRoomsOwnerHandler - список комнат, выставленных конкретным пользователем(арендодателем)
func (r *Router) GetRoomsOwnerHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respond(c, http.StatusUnauthorized, nil, UNAUTH)
		return
	}

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

	rooms, err := r.store.Room().GetRooms(context.Background(), limitInt, offsetInt, map[string]string{
		"owner_id": fmt.Sprintf("=%d", userID.(int64)),
	}, false, 0, 0, 0)

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

// GetLotsOwnerHandler - список квартир, выставленных конкретным пользователем(арендодателем)
func (r *Router) GetLotsOwnerHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respond(c, http.StatusUnauthorized, nil, UNAUTH)
		return
	}

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

	rooms, err := r.store.Lot().GetFlats(context.Background(), limitInt, offsetInt, map[string]string{
		"owner_id": fmt.Sprintf("=%d", userID.(int64)),
	}, false, nil, 0, 0, 0)

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

// GetConstructOwnerHandler - список конструкторов квартир, выставленных конкретным пользователем(арендодателем)
func (r *Router) GetConstructOwnerHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respond(c, http.StatusUnauthorized, nil, UNAUTH)
		return
	}

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

	rooms, err := r.store.Lot().GetFlats(context.Background(), limitInt, offsetInt, map[string]string{
		"owner_id": fmt.Sprintf("=%d", userID.(int64)),
	}, true, nil, 0, 0, 0)

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

func (r *Router) OptionsHandler(c *gin.Context) {
	c.Status(204)
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

func mapOfParams(c *gin.Context, fields map[string]int8) map[string]string {
	result := make(map[string]string)
	for key := range fields {
		if value := c.Query(key); value != "" {
			result[key] = value
		}
	}

	log.Println(result)
	return result
}
