package postgresstore

import (
	"github.com/Vysogota99/school/internal/auth/store"
)

// StorePSQL - реализует взаимодействие с базой данных
type StorePSQL struct {
	userRepository *UserRepository
}

// New - инициализирует Store
func New() *StorePSQL {
	return &StorePSQL{}
}

// User ...
func (s *StorePSQL) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}
