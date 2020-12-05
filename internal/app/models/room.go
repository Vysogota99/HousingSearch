package models

// Room - структура описывающая комнату в квартире
type Room struct {
	ID                    int           `json:"id,omitempty"`
	LivingPlaces          []LivingPlace `binding:"required" json:"living_place,omitempty"`
	FlatID                int           `json:"flat_id,omitempty"`
	MaxResidents          int           `binding:"required" json:"max_residents"`
	CurrNumberOfResidents int           `json:"curr_number_of_residents,omitempty"`
	Description           string        `binding:"required" json:"description,omitempty"`
	Windows               bool          `binding:"required" json:"windows"`
	Balcony               bool          `json:"balcony"`
	NumOfTables           int           `binding:"required" json:"num_of_tables"`
	NumOfChairs           int           `binding:"required" json:"num_of_chairs"`
	TV                    bool          `json:"tv"`
	Furniture             bool          `binding:"required" json:"furnniture"`
	Area                  int           `binding:"required" json:"area,omitempty"`
	AvgDeposit            float64       `json:"avg_deposit,omitempty"`
	AvgPrice              float64       `json:"avg_price,omitempty"`
	Price                 float64       `json:"room_price,omitempty"`
	Deposit               float64       `json:"room_deposit,omitempty"`
	IsVisible             bool          `json:"is_visible"`
}

// RoomWorker ...
type RoomWorker struct {
	Rooms  []Room
	FlatID int
	Error  error
}

// RoomExtended - расширенная структура Room с добавление полей структуры Lot
type RoomExtended struct {
	ID                     int           `json:"id,omitempty"`
	LivingPlaces           []LivingPlace `json:"living_place,omitempty"`
	FlatID                 int           `json:"flat_id,omitempty"`
	MaxResidents           int           `json:"max_residents"`
	CurrNumberOfResidents  int           `json:"curr_number_of_residents"`
	Description            string        `json:"description,omitempty"`
	Window                 bool          `binding:"required" json:"window"`
	Balcony                bool          `json:"balcony"`
	NumOfTables            int           `json:"num_of_tables"`
	NumOfChairs            int           `json:"num_of_chairs"`
	TV                     bool          `json:"tv"`
	Furniture              bool          `binding:"required" json:"furnniture"`
	Area                   int           `json:"area"`
	AvgDeposit             float64       `json:"avg_deposit,omitempty"`
	AvgPrice               float64       `json:"avg_price,omitempty"`
	Address                string        `json:"address,omitempty"`
	Floor                  int           `json:"floor,omitempty"`
	FloorsTotal            int           `json:"floor_total,omitempty"`
	MetroStation           string        `json:"metro,omitempty"`
	TimeToMetroByTransport int           `json:"ttmetro_transport,omitempty"`
	TimeToMetroByFoot      int           `json:"ttmetro_foot,omitempty"`
	FlatArea               int           `json:"flat_area,omitempty"`
	Long                   float64       `json:"long,omitempty"`
	Lat                    float64       `json:"lat,omitempty"`
	Price                  float64       `json:"room_price,omitempty"`
	Deposit                float64       `json:"room_deposit,omitempty"`
	IsVisible              bool          `json:"is_visible"`
}

// MapRoom - ...
var MapRoom = map[string]int8{
	"area min":                 1,
	"area max":                 1,
	"max_residents":            1,
	"curr_number_of_residents": 1,
	"windows":                  1,
	"balcony":                  1,
	"num_of_tables":            1,
	"num_of_chairs":            1,
	"tv":                       1,
	"furniture":                1,
	"deposit":                  1,
	"price max":                1,
	"price min":                1,
}
