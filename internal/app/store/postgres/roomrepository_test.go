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
	res, err := store.Room().GetRoomsForMap(context.Background(), limit)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
