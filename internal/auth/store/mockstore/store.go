package mockstore

import (
	"github.com/Vysogota99/school/internal/auth/store"
)

// StoreMock - реализует взаимодействие с базой данных
type StoreMock struct {
	userRepository store.UserRepository
}

// New - инициализирует Store
func New() *StoreMock {
	return &StoreMock{}
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
