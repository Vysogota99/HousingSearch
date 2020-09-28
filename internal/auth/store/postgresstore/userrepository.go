package postgresstore

import (
	"github.com/Vysogota99/school/pkg/authService"
)

// UserRepository - реализует функционал модели User
type UserRepository struct {
	store *StorePSQL
}

// CreateUser - создает нового пользователя в базе данных
func (u *UserRepository) CreateUser(user *authService.User) error {
	return nil
}