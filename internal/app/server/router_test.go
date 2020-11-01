package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	"github.com/Vysogota99/HousingSearch/internal/app/store/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort         = ":8081"
	connStringPostgres = "user=housing_admin password=admin dbname=housing sslmode=disable"
	secretCookieKey    = "secret"
	authServicePort    = ":3001"
)

func TestGetLotsHandler(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	router := NewRouter(serverPort, store, nil)

	ts := httptest.NewServer(router.Setup())
	defer ts.Close()

	_, err := http.Get(fmt.Sprintf("%s/api/lot", ts.URL))
	assert.NoError(t, err)
}

func TestGetLotHandler(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	router := NewRouter(serverPort, store, nil)

	ts := httptest.NewServer(router.Setup())
	defer ts.Close()

	_, err := http.Get(fmt.Sprintf("%s/api/lot/1", ts.URL))
	assert.NoError(t, err)
}

func TestPostRoomHandler(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	router := NewRouter(serverPort, store, nil)
	ts := httptest.NewServer(router.Setup())

	lp := []models.LivingPlace{
		models.LivingPlace{
			ResidentID:  1,
			Price:       30000,
			Description: "Место на двоих",
			NumOFBerth:  2,
			Deposit:     15000,
		},
	}
	room := models.Room{
		LivingPlaces:   lp,
		FlatID:         1,
		MaxResidents:   2,
		Description:    "text",
		NumOfWindows:   1,
		Balcony:        false,
		NumOfTables:    1,
		NumOfChairs:    2,
		TV:             false,
		NumOFCupboards: 1,
		Area:           30,
	}

	data, err := json.Marshal(room)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/api/rooms", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, "200", resp.StatusCode)
}

func TestGetRooms(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	router := NewRouter(serverPort, store, nil)
	ts := httptest.NewServer(router.Setup())

	_, err := http.Get(fmt.Sprintf("%s/api/rooms?balcony==true&tv==false&area=>15", ts.URL))
	assert.NoError(t, err)
}
