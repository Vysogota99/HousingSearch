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
	GetFlats(ctx context.Context, limit, offset int, filters map[string]string, orderBy []string, long, lat float64, radius int) (models.Paginations, error)
	// GetFlatsFiltered(context.Context, int, int, ...map[string]string) ([]models.Lot, error)
	GetFlatAd(context.Context, int) (*models.Lot, error)
	// GetFlatAd(context.Context, int) (*models.Lot, error)
	Create(context.Context, *models.Lot) error
}

// RoomRepository - интерфейс для работы с структурой Room
type RoomRepository interface {
	GetRooms(ctx context.Context, limit, offset int, filters map[string]string, long, lat float64, radius int) (models.PaginationsRoom, error)
	GetRoom(ctx context.Context, id int) (*models.Room, error)
	Create(context.Context, *models.Room) error
	GetLivingPlace(ctx context.Context, id int) ([]models.LivingPlace, error)
	UpdateRoom(ctx context.Context, id int, fields map[string]interface{}) error
	UpdateLivingPlace(ctx context.Context, id int, fields map[string]interface{}) error
	DeleteRoom(ctx context.Context, id int) error
	DeleteLivingPlace(ctx context.Context, id int) error
}
