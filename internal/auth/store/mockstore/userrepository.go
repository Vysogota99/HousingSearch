package mockstore

import (
	"github.com/Vysogota99/school/pkg/authService"
)

// UserRepository - реализует функционал модели User
type UserRepository struct {
	store *StoreMock
}

// CreateUser - создает нового пользователя в базе данных
func (u *UserRepository) CreateUser(user *authService.User) error {

	return nil
}

// GetUser - получение пользователя (по номеру телефона) из базы данных
func (u *UserRepository) GetUser(telephoneNumber string) (*authService.User, error) {
	return nil, nil
}
