package postgres

import (
	"context"
	"database/sql"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// RoomRepository ...
type RoomRepository struct {
	store *StorePSQL
}

// GetRoomsForMap ...
func (r *RoomRepository) GetRoomsForMap(ctx context.Context, limit int) ([]models.RoomExtended, error) {
	db, err := sql.Open("postgres", r.store.ConnString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	query := `SELECT r.flatid, r.area, r.id, r.maxresidents, r.currnumberofresidents, lp.avg_price, lp.avg_deposit,
				f.address, f.floor, f.floortotal, f.metrostation, f.timetometrobytransport, f.area,
				f.long, f.lat
				FROM (
					SELECT roomid, AVG(price) AS avg_price, AVG(deposit) AS avg_deposit
					FROM living_places
					GROUP BY roomid
				) as lp
				INNER JOIN rooms r ON r.id = lp.roomid
				INNER JOIN flats f ON f.id = r.flatid
				LIMIT $1
	`
	rows, err := tx.QueryContext(ctx, query, limit)

	rooms := []models.RoomExtended{}
	for rows.Next() {
		room := models.RoomExtended{}
		if err := rows.Scan(&room.FlatID, &room.Area, &room.ID, &room.MaxResidents, &room.CurrNumberOfResidents, &room.AvgPrice, &room.AvgDeposit,
			&room.Address, &room.Floor, &room.FloorsTotal, &room.MetroStation, &room.TimeToMetroByTransport, &room.FlatArea, &room.Long, &room.Lat,
		); err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetRoom ...
func (r *RoomRepository) GetRoom(ctx context.Context, id int) (*models.Lot, error) {
	return nil, nil
}

// Create ...
func (r *RoomRepository) Create(ctx context.Context, room *models.Room) error {
	db, err := sql.Open("postgres", r.store.ConnString)
	defer db.Close()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, `
	INSERT INTO rooms(flatid, maxresidents, description, currnumberofresidents, numofwindows,
					  balcony, numoftables, numofchairs, tv, numofcupboards, area)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
	`,
		room.FlatID, room.MaxResidents, room.Description, room.CurrNumberOfResidents, room.NumOfWindows,
		room.Balcony, room.NumOfTables, room.NumOfChairs, room.TV, room.NumOFCupboards, room.Area,
	).Scan(&room.ID)
	if err != nil {
		return err
	}

	for _, lp := range room.LivingPlaces {
		err := tx.QueryRowContext(ctx, `
								INSERT INTO living_places(roomid, price, description, numofberths, deposit)
								VALUES($1, $2, $3, $4, $5) RETURNING id 
					`,
			room.ID, lp.Price, lp.Description, lp.NumOFBerth, lp.Deposit,
		).Scan(&lp.ID)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
