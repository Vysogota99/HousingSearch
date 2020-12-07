package postgres

import (
	"context"
	"testing"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/stretchr/testify/assert"
)

const (
	STORAGE_LEVEL = 15
)

func TestCreateRoom(t *testing.T) {
	flatiD := 3
	var store store.Store = New(connString, STORAGE_LEVEL)
	lot1 := models.TestLot
	lot1.ID = flatiD
	lot1.Rooms[0].FlatID = flatiD

	err := store.Room().Create(context.Background(), &lot1.Rooms[0])
	assert.NoError(t, err)
}

func TestGetRooms(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	limit := 100
	offset := 1
	filters := map[string]string{
		"area max": "<=40",
		"area min": ">5",
		"balcony":  "=true",
	}

	_, err := store.Room().GetRooms(context.Background(), limit, offset, filters, false, 0, 0, 0)
	assert.NoError(t, err)
}
func TestGetRoomsForMap(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	limit := 100
	offset := 1
	res, err := store.Room().GetRooms(context.Background(), limit, offset, nil, false, 0, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetRoom(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	var id = 5
	r, err := store.Room().GetRoom(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestGetLP(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	var id = 1
	lpArr, err := store.Room().GetLivingPlace(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, lpArr)
}

func TestUpdateRoom(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	var id = 1
	fields := map[string]interface{}{
		"description": "",
		"area":        5,
	}
	err := store.Room().UpdateRoom(context.Background(), id, fields)

	assert.NoError(t, err)
}

func TestUpdateLP(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	var id = 5
	fields := map[string]interface{}{
		"price": 30,
	}
	err := store.Room().UpdateLivingPlace(context.Background(), id, fields)

	assert.NoError(t, err)
}

func TestDeleteRoom(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	var id = 1
	err := store.Room().DeleteRoom(context.Background(), id)

	assert.NoError(t, err)
}
