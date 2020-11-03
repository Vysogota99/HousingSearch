package postgres

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// LotRepository ...
type LotRepository struct {
	store *StorePSQL
}

// GetFlats - выводит список всех квартир(объявлений)
func (l *LotRepository) GetFlats(ctx context.Context, limit, offset int, params map[string][2]string, orderBy []string) (models.Paginations, error) {
	db, err := sql.Open("postgres", l.store.ConnString)
	result := models.Paginations{}
	result.CurrentPage = offset

	defer db.Close()
	if err != nil {
		return result, err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()

	if err != nil {
		return result, err
	}

	queryNumPages := `
					SELECT CAST (count(f.id)/$1 + 1 AS integer) as num_pages 
					FROM (
						SELECT flatid, sum(maxresidents) as maxresidents, sum(currnumberofresidents) as currnumberofresidents
						FROM rooms
						GROUP BY flatid
					) as r
					INNER JOIN flats f ON f.id = r.flatid
					WHERE f.isvisible = true
	   `
	if err := tx.QueryRowContext(ctx, queryNumPages, limit).Scan(&result.NumPages); err != nil {
		return result, err
	}
	var queryFlats string
	var queryRooms string

	if params != nil {
		queryFlats = `
					SELECT f.id, f.price, f.deposit, f.address, f.floor, f.floortotal, f.metrostation, f.timetometrobytransport, f.area,
							r.maxresidents, r.currnumberofresidents, f.long, f.lat
					FROM (
						SELECT flatid, sum(maxresidents) as maxresidents, sum(currnumberofresidents) as currnumberofresidents
						FROM rooms
						GROUP BY flatid
					) as r
					INNER JOIN flats f ON f.id = r.flatid
					WHERE f.isvisible = true `

		for key, value := range params {
			queryFlats += " AND f." + key + value[0] + value[1]
		}

		queryFlats += "ORDER BY " + orderBy[0] + " " + orderBy[1] + " LIMIT $1 OFFSET $1 * ($2 - 1)"

	} else {
		queryFlats = `
			SELECT f.id, f.price, f.deposit, f.address, f.floor, f.floortotal, f.metrostation, f.timetometrobytransport, f.area,
					r.maxresidents, r.currnumberofresidents, f.long, f.lat
			FROM (
				SELECT flatid, sum(maxresidents) as maxresidents, sum(currnumberofresidents) as currnumberofresidents
				FROM rooms
				GROUP BY flatid
			) as r
			INNER JOIN flats f ON f.id = r.flatid
			WHERE f.isvisible = true ` + "ORDER BY " + orderBy[0] + " " + orderBy[1] + " LIMIT $1 OFFSET $1 * ($2 - 1)"
	}

	log.Println(queryFlats)
	rowsFlats, err := tx.QueryContext(ctx, queryFlats, limit, offset)
	defer rowsFlats.Close()
	if err != nil {
		return result, err
	}

	i := 0
	flatsID := []interface{}{}
	lots := []models.Lot{}

	for rowsFlats.Next() {
		lot := models.Lot{}
		coord := models.Point{}
		lot.Coordinates = coord
		if err := rowsFlats.Scan(&lot.ID, &lot.Price, &lot.Deposit, &lot.Address, &lot.Floor, &lot.FloorsTotal, &lot.MetroStation, &lot.TimeToMetroByTransport, &lot.Area, &lot.TotalNumberOfResidents, &lot.CurrNumberOfResidents, &lot.Coordinates.X, &lot.Coordinates.Y); err != nil {
			return result, err
		}
		lots = append(lots, lot)
		flatsID = append(flatsID, lot.ID)
		i++
	}
	if cap(flatsID) == 0 {
		return result, sql.ErrNoRows
	}

	queryRooms = `
			SELECT r.flatid, r.id, r.price, r.deposit, r.maxresidents, r.currnumberofresidents, lp.avg_price, lp.avg_deposit
			FROM (
				SELECT roomid, AVG(price) AS avg_price, AVG(deposit) AS avg_deposit
				FROM living_places
				GROUP BY roomid
			) as lp
			INNER JOIN rooms r ON r.id = lp.roomid
			WHERE r.flatid in ($1`

	for j := 2; j <= len(flatsID); j++ {
		queryRooms += ", $" + strconv.Itoa(j)
	}
	queryRooms += ")"
	log.Println(queryRooms)

	rowsRooms, err := tx.QueryContext(ctx, queryRooms, flatsID...)
	defer rowsRooms.Close()
	if err != nil {
		return result, err
	}

	dictWithRooms := make(map[string][]models.Room)
	for rowsRooms.Next() {
		room := models.Room{}
		if err := rowsRooms.Scan(&room.FlatID, &room.ID, &room.Price, &room.Deposit, &room.MaxResidents, &room.CurrNumberOfResidents, &room.AvgPrice, &room.AvgDeposit); err != nil {
			return result, err
		}

		dictWithRooms[strconv.Itoa(room.FlatID)] = append(dictWithRooms[strconv.Itoa(room.FlatID)], room)
	}

	for i = 0; i < len(lots); i++ {
		lots[i].Rooms = dictWithRooms[strconv.Itoa(lots[i].ID)]
	}

	tx.Commit()
	result.Data = lots
	return result, nil
}

// GetFlatsFiltered - выводит список всех квартир(объявлений) с применением фильтров и сортировок
func (l *LotRepository) GetFlatsFiltered(ctx context.Context, limit, offset int, params ...map[string]string) ([]models.Lot, error) {

	return nil, nil
}

// GetFlat - выводит конкретную квартиру(объявление)
func (l *LotRepository) GetFlat(ctx context.Context, id int) (*models.Lot, error) {
	db, err := sql.Open("postgres", l.store.ConnString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	lot := models.Lot{}
	coord := models.Point{}
	lot.Coordinates = coord
	queryFlat := "SELECT * FROM flats WHERE id = $1"
	err = tx.QueryRowContext(ctx, queryFlat, id).Scan(&lot.ID, &lot.OwnerID, &lot.Address, &lot.Coordinates.X, &lot.Coordinates.Y, &lot.Price, &lot.Deposit, &lot.Description,
		&lot.TimeToMetroONFoot, &lot.TimeToMetroByTransport, &lot.MetroStation, &lot.Floor, &lot.FloorsTotal,
		&lot.Area, &lot.Repairs, &lot.Elevators, &lot.Bathroom, &lot.Refrigerator, &lot.Dishwasher, &lot.GasStove,
		&lot.ElectricStove, &lot.VacuumCleaner, &lot.Internet, &lot.Animals, &lot.Smoking, &lot.IsVisible, &lot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryRooms := "SELECT * FROM rooms WHERE flatid = $1"
	rooms := []models.Room{}
	rowsRooms, err := tx.QueryContext(ctx, queryRooms, id)
	defer rowsRooms.Close()
	if err != nil {
		return nil, err
	}

	roomsID := []interface{}{}
	for rowsRooms.Next() {
		room := models.Room{}
		if err := rowsRooms.Scan(&room.ID, &room.FlatID, &room.MaxResidents, &room.Description, &room.Price, &room.Deposit, &room.CurrNumberOfResidents, &room.NumOfWindows,
			&room.Balcony, &room.NumOfTables, &room.NumOfChairs, &room.TV, &room.NumOFCupboards, &room.Area,
		); err != nil {
			return nil, err
		}

		roomsID = append(roomsID, room.ID)
		rooms = append(rooms, room)
	}

	queryLivingPlaces := "SELECT * FROM living_places WHERE roomid IN ($1"
	for i := 2; i <= cap(roomsID); i++ {
		queryLivingPlaces += ", $" + strconv.Itoa(i)
	}
	queryLivingPlaces += ")"

	rowsLivingPlaces, err := tx.QueryContext(ctx, queryLivingPlaces, roomsID...)
	defer rowsLivingPlaces.Close()
	if err != nil {
		return nil, err
	}

	dictWithLPlaces := make(map[string][]models.LivingPlace)
	for rowsLivingPlaces.Next() {
		lp := models.LivingPlace{}
		var residentID sql.NullInt64
		if err := rowsLivingPlaces.Scan(&lp.ID, &lp.RoomID, &residentID, &lp.Price, &lp.Description, &lp.NumOFBerth, &lp.Description); err != nil {
			return nil, err
		}

		if residentID.Valid {
			lp.ResidentID = int(residentID.Int64)
		}

		dictWithLPlaces[strconv.Itoa(lp.RoomID)] = append(dictWithLPlaces[strconv.Itoa(lp.RoomID)], lp)
	}

	for i := 0; i < cap(rooms); i++ {
		rooms[i].LivingPlaces = dictWithLPlaces[strconv.Itoa(rooms[i].ID)]
	}

	lot.Rooms = rooms
	tx.Commit()
	return &lot, nil
}

// Create - создает объявление
func (l *LotRepository) Create(ctx context.Context, lot *models.Lot) error {
	db, err := sql.Open("postgres", l.store.ConnString)
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
									INSERT INTO flats(ownerID, address, long, lat, description, 
													  timeToMetroOnFoot, timeToMetroByTransport, metrostation,
													  floor, floorTotal, area, repair, elevator,
													  bathroom, refrigerator, dishwasher, gasStove,
													  electricStove, vacuumcleaner, internet,
													  animals, smoking)
									VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) 
									RETURNING id
								`,
		lot.OwnerID, lot.Address, lot.Coordinates.X, lot.Coordinates.Y, lot.Description, lot.TimeToMetroONFoot, lot.TimeToMetroByTransport,
		lot.MetroStation, lot.Floor, lot.FloorsTotal, lot.Area, lot.Repairs, lot.Elevators, lot.Bathroom, lot.Refrigerator, lot.Dishwasher,
		lot.GasStove, lot.ElectricStove, lot.VacuumCleaner, lot.Internet, lot.Animals, lot.Smoking,
	).Scan(&lot.ID)
	if err != nil {
		return err
	}

	for i := 0; i < len(lot.Rooms); i++ {
		err = tx.QueryRowContext(ctx, `
				INSERT INTO rooms(flatid, maxresidents, description, currnumberofresidents, numofwindows,
								  balcony, numoftables, numofchairs, tv, numofcupboards, area)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING id
		`,
			lot.ID, lot.Rooms[i].MaxResidents, lot.Rooms[i].Description, lot.Rooms[i].CurrNumberOfResidents, lot.Rooms[i].NumOfWindows,
			lot.Rooms[i].Balcony, lot.Rooms[i].NumOfTables, lot.Rooms[i].NumOfChairs, lot.Rooms[i].TV, lot.Rooms[i].NumOFCupboards, lot.Rooms[i].Area,
		).Scan(&lot.Rooms[i].ID)
		if err != nil {
			return err
		}

		for j := 0; j < len(lot.Rooms[i].LivingPlaces); j++ {
			err := tx.QueryRowContext(ctx, `
								INSERT INTO living_places(roomid, description, numofberths)
								VALUES($1, $2, $3) RETURNING id 
					`,
				lot.Rooms[i].ID, lot.Rooms[i].LivingPlaces[j].Description, lot.Rooms[i].LivingPlaces[j].NumOFBerth,
			).Scan(&lot.Rooms[i].LivingPlaces[j].ID)
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()
	return nil
}
