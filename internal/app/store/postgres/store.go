package postgres

import (
	"github.com/Vysogota99/HousingSearch/internal/app/store"
	_ "github.com/lib/pq"
)

// StorePSQL - реализует взаимодействие с базой данных
type StorePSQL struct {
	ConnString     string
	lotRepository  store.LotRepository
	roomRepository store.RoomRepository
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

// Room ...
func (s *StorePSQL) Room() store.RoomRepository {
	if s.roomRepository == nil {
		s.roomRepository = &RoomRepository{
			store: s,
		}
	}

	return s.roomRepository
}
