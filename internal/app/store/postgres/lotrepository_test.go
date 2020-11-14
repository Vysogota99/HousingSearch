package postgres

import (
	"context"
	"testing"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/stretchr/testify/assert"
)

const (
	connString = "user=housing_admin password=admin dbname=housing sslmode=disable"
)

func TestCreateLot(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	lot1 := models.TestLot
	lot2 := models.TestLot

	lot1.Repairs = 2
	lot1.Smoking = false
	lot1.Address = "Россия, Москва, Нахимовский проезд, 23к1, кв21"
	lot1.Rooms[0].Area = 15
	lot1.Rooms[1].Area = 10
	lot1.Area = 35
	lot1.IsConstructor = false

	lot2.Address = "Россия, Москва, Коломенский проезд, 23к1, кв15"
	lot2.Rooms[0].LivingPlaces[0].Price = 30000
	lot2.Rooms[0].LivingPlaces[1].Price = 20000
	lot2.Area = 120
	err := store.Lot().Create(context.Background(), &lot1)
	assert.NoError(t, err)

	err = store.Lot().Create(context.Background(), &lot2)
	assert.NoError(t, err)
}

func TestGetLots(t *testing.T) {
	// var store store.Store = New(connString, STORAGE_LEVEL)
	// orderBy := []string{"created_at", "desc"}
	// lots, err := store.Lot().GetFlats(context.Background(), 3, 1, nil, orderBy)
	// assert.NoError(t, err)
	// assert.NotNil(t, lots)
}

func TestGetLotsFiltered(t *testing.T) {
	// var store store.Store = New(connString, STORAGE_LEVEL)

	// filters := make(map[string][2]string)

	// keysArea := [2]string{">=", "15"}
	// orderBy := []string{"area", "asc"}
	// filters["area"] = keysArea

	// lots, err := store.Lot().GetFlats(context.Background(), 10, 1, filters, orderBy)
	// assert.NoError(t, err)
	// assert.NotNil(t, lots)
}

func TestGetFlatAd(t *testing.T) {
	var store store.Store = New(connString, STORAGE_LEVEL)
	lot, err := store.Lot().GetFlatAd(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, lot)
}
