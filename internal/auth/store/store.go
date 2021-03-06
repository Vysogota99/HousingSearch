package store

import "github.com/Vysogota99/HousingSearch/pkg/authService"

// Store - ...
type Store interface {
	User() UserRepository
}

// UserRepository ...
type UserRepository interface {
	CreateUser(user *authService.User) error
	GetUser(telephoneNumber string) (*authService.User, error)
}
