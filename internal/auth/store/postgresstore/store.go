package postgresstore

import (
	"github.com/Vysogota99/HousingSearch/internal/auth/store"
)

// StorePSQL - реализует взаимодействие с базой данных
type StorePSQL struct {
	userRepository *UserRepository
	ConnString     string
}

// New - инициализирует Store
func New(connString string) *StorePSQL {
	return &StorePSQL{
		ConnString: connString,
	}
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
