package mockstore

import (
	"github.com/Vysogota99/HousingSearch/internal/auth/store"
)

// StoreMock - реализует взаимодействие с базой данных
type StoreMock struct {
	userRepository store.UserRepository
	ConnString     string
}

// New - инициализирует Store
func New(connString string) *StoreMock {
	return &StoreMock{
		ConnString: connString,
	}
}

// User ...
func (s *StoreMock) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}
