package store

import "github.com/Vysogota99/school/internal/auth/models"

// Store - ...
type Store interface {
	User() UserRepository
}

// UserRepository ...
type UserRepository interface {
	CreateUser(user *models.User) error
}
