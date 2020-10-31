package store

import (
	"context"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// Store - ...
type Store interface {
	Lot() LotRepository
	Room() RoomRepository
}

// LotRepository - лот с жильем. При обращении к лоту, можно получить выборку
// квартир, комнат, спальных мест с применением различных фильтров
type LotRepository interface {
	GetFlats(ctx context.Context, limit, offset int, params map[string][2]string, orderBy []string) (models.Paginations, error)
	// GetFlatsFiltered(context.Context, int, int, ...map[string]string) ([]models.Lot, error)
	GetFlat(context.Context, int) (*models.Lot, error)
	Create(context.Context, *models.Lot) error
}

// RoomRepository - интерфейс для работы с структурой Room
type RoomRepository interface {
	GetRoomsForMap(ctx context.Context, limit int) ([]models.RoomExtended, error)
	GetRoom(ctx context.Context, id int) (*models.Lot, error)
	Create(context.Context, *models.Room) error
}
