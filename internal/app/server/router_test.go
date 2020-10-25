package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetLot(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	router := NewRouter(serverPort, store, nil)

	ts := httptest.NewServer(router.Setup())
	defer ts.Close()

	_, err := http.Get(fmt.Sprintf("%s/api/lot/1", ts.URL))
	assert.NoError(t, err)
}
