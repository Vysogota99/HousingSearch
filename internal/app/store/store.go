package store

import (
	"context"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// Store - ...
type Store interface {
	Lot() LotRepository
}

// LotRepository - лот с жильем. При обращении к лоту, можно получить выборку
// квартир, комнат, спальных мест с применением различных фильтров
type LotRepository interface {
	GetFlats(ctx context.Context, limit, offset int, params map[string][2]string, orderBy [2]string) (models.Paginations, error)
	// GetFlatsFiltered(context.Context, int, int, ...map[string]string) ([]models.Lot, error)
	GetFlat(context.Context, int) (*models.Lot, error)
	Create(context.Context, *models.Lot) error
}

// RoomRepository - интерфейс для работы с структурой Room
type RoomRepository interface {
	GetRooms() (context.Context, []models.Room, error)
	GetRoomsFiltered(context.Context, ...map[string]string) ([]models.Room, error)
	GetRoom(int) (context.Context, *models.Lot, error)
	Create(context.Context, *models.Room) error
}
