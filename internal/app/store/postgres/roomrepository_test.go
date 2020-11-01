package postgres

import (
	"context"
	"testing"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {
	var store store.Store = New(connString)
	lot1 := models.TestLot
	lot1.ID = 1
	lot1.Rooms[0].FlatID = 1

	err := store.Room().Create(context.Background(), &lot1.Rooms[0])
	assert.NoError(t, err)
}

func TestGetRoomsForMap(t *testing.T) {
	var store store.Store = New(connString)
	limit := 100
	offset := 1
	res, err := store.Room().GetRooms(context.Background(), limit, offset, nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGetRoom(t *testing.T) {
	var store store.Store = New(connString)
	var id = 1
	r, err := store.Room().GetRoom(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestGetLP(t *testing.T) {
	var store store.Store = New(connString)
	var id = 1
	lpArr, err := store.Room().GetLivingPlace(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, lpArr)
}

func TestUpdateRoom(t *testing.T) {
	var store store.Store = New(connString)
	var id = 1
	// fields := map[string]interface{}{
	// 	"area": 30,
	// 	"tv":   true,
	// }
	err := store.Room().UpdateRoom(context.Background(), id, nil)

	assert.NoError(t, err)
}

func TestUpdateLP(t *testing.T) {
	var store store.Store = New(connString)
	var id = 5
	fields := map[string]interface{}{
		"price": 30,
	}
	err := store.Room().UpdateLivingPlace(context.Background(), id, fields)

	assert.NoError(t, err)
}

func TestDeleteRoom(t *testing.T) {
	var store store.Store = New(connString)
	var id = 1
	err := store.Room().DeleteRoom(context.Background(), id)

	assert.NoError(t, err)
}
