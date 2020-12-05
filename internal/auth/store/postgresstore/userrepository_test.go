package postgresstore_test

import (
	"testing"

	"github.com/Vysogota99/HousingSearch/internal/auth/store/postgresstore"
	"github.com/Vysogota99/HousingSearch/pkg/authService"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	connstring := "user=housing_admin password=admin dbname=housing sslmode=disable"
	store := postgresstore.New(connstring)

	user := &authService.User{}
	user.PassName = "Ivan"
	user.PassLastName = "lapshin"
	user.PassSex = "male"
	user.PassDateOfBirth = "03.09.1999"
	user.Password = "qwert"
	user.Role = 1
	user.TelephoneNumber = "+79037658681"

	err := store.User().CreateUser(user)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	connstring := "user=auth_user password=admin dbname=auth_user sslmode=disable"
	store := postgresstore.New(connstring)
	user, err := store.User().GetUser("89037658681")
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
