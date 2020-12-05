package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
	"github.com/golang/geo/s2"
	_ "github.com/lib/pq"
)

// LotRepository ...
type LotRepository struct {
	store        *StorePSQL
	storageLevel int
}

const (
	STOP = "Вызвана остановка"
)

// GetFlats - постранично выводит список квартир(объявлений)
func (l *LotRepository) GetFlats(ctx context.Context, limit, offset int, filters map[string]string, isConstruct bool, orderBy []string, long, lat float64, radius int) (models.Paginations, error) {
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

	// Поиск квартиры в радиусе
	condition := ""
	if long != 0 || lat != 0 || radius != 0 {
		condition = "WHERE cell_id IN (%s) AND "
		cu := searchCellUnion(long, lat, radius, l.store.storageLevel)
		cellIDS := make([]string, 0)
		for _, item := range cu {
			cellIDS = append(cellIDS, strconv.Itoa(int(item)))
		}

		cellIDSStr := strings.Join(cellIDS, ",")
		condition = fmt.Sprintf(condition, cellIDSStr)
	}
	// Дополнительные фильтры
	if filters != nil && len(filters) > 0 {
		if condition == "" {
			condition = "WHERE "
		}

		for key, value := range filters {
			condition += key + value + " AND "
		}
	}

	if condition == "" {
		condition = fmt.Sprintf("WHERE is_visible = true AND is_constructor = %t", isConstruct)
	} else {
		condition += fmt.Sprintf("is_visible = true AND is_constructor = %t", isConstruct)
	}
	// Конец фильтров

	queryNumPages := `SELECT CAST (count(id)/$1 + 1 AS integer) as num_pages 
				      FROM flats
					  %s`
	queryNumPages = fmt.Sprintf(queryNumPages, condition)
	if err := tx.QueryRowContext(ctx, queryNumPages, limit).Scan(&result.NumPages); err != nil {
		return result, err
	}

	queryFlats := ` SELECT id, owner_id, address, long, long, price, deposit, description, time_to_metro_on_foot,
					time_to_metro_by_transport, metro_station, floor, floor_total, area, repair, pass_elevator,
					service_elevator, kitchen, microwave_oven, bathroom, refrigerator, dishwasher, stove, vacuum_cleaner,
					dryer, internet, animals, smoking, heating, is_visible, is_constructor, created_at, updated_at
					FROM flats
					%s
					LIMIT $1 OFFSET $1 * ($2 - 1)`
	queryFlats = fmt.Sprintf(queryFlats, condition)
	// log.Println(queryFlats)

	rows, err := tx.QueryContext(ctx, queryFlats, limit, offset)
	if err != nil {
		return result, err
	}

	flatsIdsSlice := make([]string, 0)
	for rows.Next() {
		lot := models.Lot{}
		lot.Coordinates = models.Point{}
		if err := rows.Scan(&lot.ID, &lot.OwnerID, &lot.Address, &lot.Coordinates.X, &lot.Coordinates.Y, &lot.Price, &lot.Deposit, &lot.Description, &lot.TimeToMetroONFoot, &lot.TimeToMetroByTransport,
			&lot.MetroStation, &lot.Floor, &lot.FloorsTotal, &lot.Area, &lot.Repairs, &lot.PassElevator, &lot.ServiceElevator, &lot.Kitchen, &lot.MicrowaveOven, &lot.Bathroom, &lot.Refrigerator, &lot.Dishwasher,
			&lot.Stove, &lot.VacuumCleaner, &lot.Dryer, &lot.Internet, &lot.Animals, &lot.Smoking, &lot.Heating, &lot.IsVisible, &lot.IsConstructor, &lot.CreatedAt, &lot.UpdatedAt,
		); err != nil {
			return result, err
		}

		flatsIdsSlice = append(flatsIdsSlice, strconv.Itoa(lot.ID))
		result.Data = append(result.Data, lot)
	}

	if len(flatsIdsSlice) == 0 {
		return result, nil
	}

	queryRooms := `SELECT id, flat_id, max_residents, description, price, deposit, curr_number_of_residents,
				   balcony, num_of_tables, num_of_chairs, tv, furniture, area, windows
				   FROM rooms WHERE ID IN (%s)`
	queryRooms = fmt.Sprintf(queryRooms, strings.Join(flatsIdsSlice, ","))
	rows, err = tx.QueryContext(ctx, queryRooms)
	if err != nil {
		return result, err
	}

	dictRooms := map[int][]models.Room{}

	for rows.Next() {
		room := models.Room{}
		if err := rows.Scan(&room.ID, &room.FlatID, &room.MaxResidents, &room.Description, &room.Price, &room.Deposit, &room.CurrNumberOfResidents,
			&room.Balcony, &room.NumOfTables, &room.NumOfChairs, &room.TV, &room.Furniture, &room.Area, &room.Windows); err != nil {
			return result, err
		}

		dictRooms[room.FlatID] = append(dictRooms[room.FlatID], room)
	}

	for index, lot := range result.Data {
		result.Data[index].Rooms = dictRooms[lot.ID]
	}
	return result, nil
}

// GetFlatAd - выводит конкретную квартиру(объявление)
func (l *LotRepository) GetFlatAd(ctx context.Context, id int, isConstructor bool) (*models.Lot, error) {
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
	queryFlat := `SELECT id, owner_id, address, long, long, price, deposit, description, time_to_metro_on_foot,
				  time_to_metro_by_transport, metro_station, floor, floor_total, area, repair, pass_elevator,
				  service_elevator, kitchen, microwave_oven, bathroom, refrigerator, dishwasher, stove, vacuum_cleaner,
				  dryer, internet, animals, smoking, heating, is_visible, is_constructor, created_at, updated_at
				  FROM flats WHERE id = $1 AND is_visible = true AND is_constructor = %t`

	queryFlat = fmt.Sprintf(queryFlat, isConstructor)
	err = tx.QueryRowContext(ctx, queryFlat, id).Scan(&lot.ID, &lot.OwnerID, &lot.Address, &lot.Coordinates.X, &lot.Coordinates.Y, &lot.Price, &lot.Deposit, &lot.Description,
		&lot.TimeToMetroONFoot, &lot.TimeToMetroByTransport, &lot.MetroStation, &lot.Floor, &lot.FloorsTotal,
		&lot.Area, &lot.Repairs, &lot.PassElevator, &lot.ServiceElevator, &lot.Kitchen, &lot.MicrowaveOven, &lot.Bathroom, &lot.Refrigerator, &lot.Dishwasher, &lot.Stove,
		&lot.VacuumCleaner, &lot.Dryer, &lot.Internet, &lot.Animals, &lot.Smoking, &lot.Heating, &lot.IsVisible, &lot.IsConstructor, &lot.CreatedAt, &lot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryRooms := `SELECT id, flat_id, max_residents, description, price,
				   deposit, curr_number_of_residents, balcony, num_of_tables,
				   num_of_chairs, tv, furniture, area, windows 
				   FROM rooms WHERE flat_id = $1`
	rooms := []models.Room{}
	rowsRooms, err := tx.QueryContext(ctx, queryRooms, id)
	defer rowsRooms.Close()
	if err != nil {
		return nil, err
	}

	roomsID := []interface{}{}
	for rowsRooms.Next() {
		room := models.Room{}
		if err := rowsRooms.Scan(&room.ID, &room.FlatID, &room.MaxResidents, &room.Description, &room.Price, &room.Deposit, &room.CurrNumberOfResidents,
			&room.Balcony, &room.NumOfTables, &room.NumOfChairs, &room.TV, &room.Furniture, &room.Area, &room.Windows,
		); err != nil {
			return nil, err
		}

		roomsID = append(roomsID, room.ID)
		rooms = append(rooms, room)
	}

	queryLivingPlaces := "SELECT id, roomid, residentid, price, description, numofberths, deposit FROM living_places WHERE roomid IN ($1"
	for i := 2; i <= len(roomsID); i++ {
		queryLivingPlaces += ", $" + strconv.Itoa(i)
	}
	queryLivingPlaces += ")"

	// log.Println(queryLivingPlaces)
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

	for i := 0; i < len(rooms); i++ {
		rooms[i].LivingPlaces = dictWithLPlaces[strconv.Itoa(rooms[i].ID)]
	}

	lot.Rooms = rooms
	tx.Commit()
	return &lot, nil
}

// Create - создает квартиру в конструкторе
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

	// нахождение номера ячейки с заданным уровнем
	latlong := s2.LatLngFromDegrees(lot.Coordinates.Y, lot.Coordinates.X)
	cellID := s2.CellIDFromLatLng(latlong)
	lot.Coordinates.CellID = uint64(cellID.Parent(l.storageLevel))

	err = tx.QueryRowContext(ctx, `
									INSERT INTO flats(owner_id, address, long, lat, cell_id,
													  time_to_metro_on_foot, time_to_metro_by_transport, metro_station,
													  floor, floor_total, area, repair, pass_elevator, service_elevator,
													  kitchen, microwave_oven, bathroom, refrigerator, dishwasher, stove,
													   vacuum_cleaner, dryer, internet, animals, smoking, heating)
									VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26) 
									RETURNING id
								`,
		lot.OwnerID, lot.Address, lot.Coordinates.X, lot.Coordinates.Y, lot.Coordinates.CellID, lot.TimeToMetroONFoot, lot.TimeToMetroByTransport,
		lot.MetroStation, lot.Floor, lot.FloorsTotal, lot.Area, lot.Repairs, lot.PassElevator, lot.ServiceElevator, lot.Kitchen, lot.MicrowaveOven, lot.Bathroom, lot.Refrigerator, lot.Dishwasher,
		lot.Stove, lot.VacuumCleaner, lot.Dryer, lot.Internet, lot.Animals, lot.Smoking, lot.Heating,
	).Scan(&lot.ID)
	if err != nil {
		return err
	}

	for i := 0; i < len(lot.Rooms); i++ {
		err = tx.QueryRowContext(ctx, `
				INSERT INTO rooms(flat_id, max_residents, curr_number_of_residents, windows,
								  balcony, num_of_tables, num_of_chairs, tv, furniture, area)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING id
		`,
			lot.ID, lot.Rooms[i].MaxResidents, lot.Rooms[i].CurrNumberOfResidents, lot.Rooms[i].Windows,
			lot.Rooms[i].Balcony, lot.Rooms[i].NumOfTables, lot.Rooms[i].NumOfChairs, lot.Rooms[i].TV, lot.Rooms[i].Furniture, lot.Rooms[i].Area,
		).Scan(&lot.Rooms[i].ID)
		if err != nil {
			return err
		}

		for j := 0; j < len(lot.Rooms[i].LivingPlaces); j++ {
			err := tx.QueryRowContext(ctx, `
								INSERT INTO living_places(roomid, numofberths)
								VALUES($1, $2) RETURNING id 
					`,
				lot.Rooms[i].ID, lot.Rooms[i].LivingPlaces[j].NumOFBerth,
			).Scan(&lot.Rooms[i].LivingPlaces[j].ID)
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()
	return nil
}
