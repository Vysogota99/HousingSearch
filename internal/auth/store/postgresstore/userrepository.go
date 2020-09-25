package postgresstore

import (
	"github.com/Vysogota99/school/internal/auth/models"
)

// UserRepository - реализует функционал модели User
type UserRepository struct {
	store *StorePSQL
}

// CreateUser - создает нового пользователя в базе данных
func (u *UserRepository) CreateUser(user *models.User) error {
	return nil
}
