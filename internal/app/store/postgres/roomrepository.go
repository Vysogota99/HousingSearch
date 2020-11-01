package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/Vysogota99/HousingSearch/internal/app/models"
)

// RoomRepository ...
type RoomRepository struct {
	store *StorePSQL
}

// GetRooms ...
func (r *RoomRepository) GetRooms(ctx context.Context, limit, offset int, filters map[string]string) ([]models.RoomExtended, error) {
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

	mapRoom := map[string]int8{
		"area":                  1,
		"maxresidents":          1,
		"currnumberofresidents": 1,
		"numofwindows":          1,
		"balcony":               1,
		"numoftables":           1,
		"numofchairs":           1,
		"tv":                    1,
		"numofcupboards":        1,
	}
	mapLp := map[string]int8{
		"avg_price":   1,
		"avg_deposit": 1,
	}
	mapFlat := map[string]int8{
		"floor":                  1,
		"floortotal":             1,
		"metrostation":           1,
		"timetometrobytransport": 1,
		"timetometroonfoot":      1,
		"long":                   1,
		"lat":                    1,
		"repair":                 1,
		"elevator":               1,
		"bathroom":               1,
		"refrigerator":           1,
		"dishwasher":             1,
		"gasstove":               1,
		"electricstove":          1,
		"vacuumcleaner":          1,
		"internet":               1,
		"animals":                1,
		"smoking":                1,
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
				%s
				LIMIT $1 OFFSET $1 * ($2 - 1)
	`
	condition := "WHERE "
	if filters != nil && len(filters) > 0 {
		for key, value := range filters {
			table := ""
			if _, ok := mapRoom[key]; ok {
				table = "r."
			}
			if _, ok := mapFlat[key]; ok {
				table = "f."
			}
			if _, ok := mapLp[key]; ok {
				table = "lp."
			}
			condition += table + key + value + " AND "
		}

		condition = condition[0 : len(condition)-4]
	}
	query = fmt.Sprintf(query, condition)
	log.Println(query)

	rows, err := tx.QueryContext(ctx, query, limit, offset)

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

	tx.Commit()
	return rooms, nil
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
	err = tx.QueryRow("SELECT * FROM rooms WHERE ID = $1", id).Scan(&room.ID, &room.FlatID, &room.MaxResidents, &room.Description, &room.CurrNumberOfResidents, &room.NumOfWindows, &room.Balcony, &room.NumOfTables, &room.NumOfChairs, &room.TV, &room.NumOFCupboards, &room.Area)
	if err != nil {
		return nil, err
	}

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
