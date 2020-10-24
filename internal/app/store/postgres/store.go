package postgres

import (
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	_ "github.com/lib/pq"
)

// StorePSQL - реализует взаимодействие с базой данных
type StorePSQL struct {
	ConnString    string
	lotRepository store.LotRepository
}

// New - инициализирует Store
func New(connString string) *StorePSQL {
	return &StorePSQL{
		ConnString: connString,
	}
}

// Lot ...
func (s *StorePSQL) Lot() store.LotRepository {
	if s.lotRepository == nil {
		s.lotRepository = &LotRepository{
			store: s,
		}
	}

	return s.lotRepository
}
