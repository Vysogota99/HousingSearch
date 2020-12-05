package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// RoomRepository ...
type RoomRepository struct {
	store *StorePSQL
}

// GetRooms ...
func (r *RoomRepository) GetRooms(ctx context.Context, limit, offset int, filters map[string]string, isConstruct bool, long, lat float64, radius int) (models.PaginationsRoom, error) {
	result := models.PaginationsRoom{}
	result.CurrentPage = offset

	db, err := sql.Open("postgres", r.store.ConnString)
	defer db.Close()
	if err != nil {
		return result, err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return result, err
	}

	// Поиск комнаты в радиусе
	condition := ""
	if long != 0 || lat != 0 || radius != 0 {
		condition = "WHERE f.cell_id IN (%s) AND "
		cu := searchCellUnion(long, lat, radius, r.store.storageLevel)
		cellIDS := make([]string, 0)
		for _, item := range cu {
			cellIDS = append(cellIDS, strconv.Itoa(int(item)))
		}

		cellIDSStr := strings.Join(cellIDS, ",")
		condition = fmt.Sprintf(condition, cellIDSStr)
	}

	// Дополнительные фильтры
	roomFields := make([]string, 0)
	if filters != nil && len(filters) > 0 {
		if condition == "" {
			condition = "WHERE "
		}

		for key, value := range filters {
			table := ""
			field := strings.Split(key, " ")
			if _, ok := models.MapRoom[key]; ok {
				table = "r."
				roomFields = append(roomFields, field[0])
			}
			if _, ok := models.MapFlat[key]; ok {
				table = "f."
			}
			if _, ok := models.MapLp[key]; ok {
				table = "lp."
			}
			condition += table + field[0] + value + " AND "
		}
	}

	if condition == "" {
		condition = fmt.Sprintf("WHERE f.is_visible = true AND r.is_visible = true AND f.is_constructor = %t", isConstruct)
	} else {
		condition += fmt.Sprintf("f.is_visible = true AND r.is_visible = true AND f.is_constructor = %t", isConstruct)
	}

	// roomFieldsStr := strings.Join(roomFields, ",")
	// Конец фильтров

	queryNumPages := `
					SELECT CAST (count(r.id)/ $1 + 1 AS integer) as num_pages
					FROM rooms AS r
					INNER JOIN flats f ON f.id = r.flat_id %s
	   `
	queryNumPages = fmt.Sprintf(queryNumPages, condition)
	log.Println(queryNumPages)
	if err := tx.QueryRowContext(ctx, queryNumPages, limit).Scan(&result.NumPages); err != nil {
		return result, err
	}

	query := `SELECT r.price, r.deposit, r.flat_id, r.area, r.id, r.max_residents, r.curr_number_of_residents, lp.avg_price, lp.avg_deposit,
				f.address, f.floor, f.floor_total, f.metro_station, f.time_to_metro_by_transport, f.time_to_metro_on_foot, f.area,
				f.long, f.lat
				FROM (
					SELECT roomid, AVG(price) AS avg_price, AVG(deposit) AS avg_deposit
					FROM living_places
					GROUP BY roomid
				) as lp
				INNER JOIN rooms r ON r.id = lp.roomid
				INNER JOIN flats f ON f.id = r.flat_id 
				%s
				LIMIT $1 OFFSET $1 * ($2 - 1)
	`
	query = fmt.Sprintf(query, condition)
	log.Println(query)

	rows, err := tx.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return result, err
	}

	rooms := []models.RoomExtended{}
	for rows.Next() {
		room := models.RoomExtended{}
		if err := rows.Scan(&room.Price, &room.Deposit, &room.FlatID, &room.Area, &room.ID, &room.MaxResidents, &room.CurrNumberOfResidents, &room.AvgPrice, &room.AvgDeposit,
			&room.Address, &room.Floor, &room.FloorsTotal, &room.MetroStation, &room.TimeToMetroByTransport, &room.TimeToMetroByFoot, &room.FlatArea, &room.Long, &room.Lat,
		); err != nil {
			return result, err
		}

		rooms = append(rooms, room)
	}

	tx.Commit()
	result.Data = rooms
	return result, nil
}

// GetRoomsAround ...
func (r *RoomRepository) GetRoomsAround(ctx context.Context, filters map[string]string, long, lat float64, radius int) (models.PaginationsRoom, error) {
	result := models.PaginationsRoom{}
	cu := searchCellUnion(long, lat, radius, r.store.storageLevel)
	cellIDS := make([]string, 0)
	for _, item := range cu {
		cellIDS = append(cellIDS, strconv.Itoa(int(item)))
	}

	cellIDSStr := strings.Join(cellIDS, ",")
	log.Print(cellIDSStr)
	return result, nil
}

// GetRoom ...
func (r *RoomRepository) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	db, err := sql.Open("postgres", r.store.ConnString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	room := models.Room{}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow(`SELECT id, flat_id, max_residents, description, price, deposit,
					   curr_number_of_residents, balcony, num_of_tables, num_of_chairs,
					   tv, furniture, area, windows 
					   FROM rooms WHERE ID = $1`, id).Scan(&room.ID, &room.FlatID, &room.MaxResidents,
		&room.Description, &room.Price, &room.Deposit, &room.CurrNumberOfResidents,
		&room.Balcony, &room.NumOfTables, &room.NumOfChairs, &room.TV,
		&room.Furniture, &room.Area, &room.Windows)
	if err != nil {
		return nil, err
	}

	lp, err := r.GetLivingPlace(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	room.LivingPlaces = lp
	tx.Commit()
	return &room, nil
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
	INSERT INTO rooms(flat_id, max_residents, description, curr_number_of_residents, windows,
					  balcony, num_of_tables, num_of_chairs, tv, furniture, area)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
	`,
		room.FlatID, room.MaxResidents, room.Description, room.CurrNumberOfResidents, room.Windows,
		room.Balcony, room.NumOfTables, room.NumOfChairs, room.TV, room.Furniture, room.Area,
	).Scan(&room.ID)
	if err != nil {
		return err
	}

	for _, lp := range room.LivingPlaces {
		err := tx.QueryRowContext(ctx, `
								INSERT INTO living_places(roomid, description, numofberths)
								VALUES($1, $2, $3) RETURNING id 
					`,
			room.ID, lp.Description, lp.NumOFBerth,
		).Scan(&lp.ID)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

// GetLivingPlace - возвращает массив спальных мест для конкретной комнаты
func (r *RoomRepository) GetLivingPlace(ctx context.Context, id int) ([]models.LivingPlace, error) {
	db, err := sql.Open("postgres", r.store.ConnString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	lpArray := []models.LivingPlace{}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT id, residentid, price, description, numofberths, deposit FROM living_places WHERE roomid = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lp := models.LivingPlace{}
		var resID sql.NullInt64
		if err := rows.Scan(&lp.ID, &resID, &lp.Price, &lp.Description, &lp.NumOFBerth, &lp.Deposit); err != nil {
			return nil, err
		}

		if resID.Valid {
			lp.ResidentID = int(resID.Int64)
		} else {
			lp.ResidentID = 0
		}

		lpArray = append(lpArray, lp)
	}

	tx.Commit()
	return lpArray, nil
}

// UpdateRoom - обновляет поля у определенной комнаты
func (r *RoomRepository) UpdateRoom(ctx context.Context, id int, fields map[string]interface{}) error {
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

	if len(fields) == 0 {
		return nil
	}

	query := "UPDATE rooms SET %s WHERE id = $1"
	data := ""
	for key, value := range fields {
		switch value.(type) {
		case int:
			value = strconv.Itoa(value.(int))
		case string:
			value = value.(string)
		case bool:
			value = strconv.FormatBool(value.(bool))
		case float32:
			value = fmt.Sprintf("%f", value.(float32))
		case float64:
			value = fmt.Sprintf("%f", value.(float64))
		}
		data += key + " = " + value.(string) + ", "
	}
	data = data[0 : len(data)-2]
	query = fmt.Sprintf(query, data)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateLivingPlace - обновляет поля у определенной комнаты
func (r *RoomRepository) UpdateLivingPlace(ctx context.Context, id int, fields map[string]interface{}) error {
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

	if len(fields) == 0 {
		return nil
	}

	query := "UPDATE living_places SET %s WHERE id = $1"
	data := ""
	for key, value := range fields {
		switch value.(type) {
		case int:
			value = strconv.Itoa(value.(int))
		case string:
			value = value.(string)
		case bool:
			value = strconv.FormatBool(value.(bool))
		case float32:
			value = fmt.Sprintf("%f", value.(float32))
		case float64:
			value = fmt.Sprintf("%f", value.(float64))
		}
		data += key + " = " + value.(string) + ", "
	}
	data = data[0 : len(data)-2]
	query = fmt.Sprintf(query, data)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// DeleteRoom ...
func (r *RoomRepository) DeleteRoom(ctx context.Context, id int) error {
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

	_, err = tx.ExecContext(ctx, "DELETE FROM rooms WHERE id = $1", id)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// DeleteLivingPlace ...
func (r *RoomRepository) DeleteLivingPlace(ctx context.Context, id int) error {
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

	_, err = tx.ExecContext(ctx, "DELETE FROM living_places WHERE id = $1", id)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
