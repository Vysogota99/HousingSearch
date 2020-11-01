package models

// Room - структура описывающая комнату в квартире
type Room struct {
	ID                    int           `json:"id,omitempty"`
	LivingPlaces          []LivingPlace `binding:"required" json:"living_place,omitempty"`
	FlatID                int           `json:"flat_id,omitempty"`
	MaxResidents          int           `binding:"required" json:"max_residents"`
	CurrNumberOfResidents int           `json:"curr_number_of_residents,omitempty"`
	Description           string        `binding:"required" json:"description,omitempty"`
	NumOfWindows          int           `binding:"required" json:"num_of_windows"`
	Balcony               bool          `json:"balcony"`
	NumOfTables           int           `binding:"required" json:"num_of_tables"`
	NumOfChairs           int           `binding:"required" json:"num_of_chairs"`
	TV                    bool          `json:"tv"`
	NumOFCupboards        int           `binding:"required" json:"num_of_cupboards"`
	Area                  int           `binding:"required" json:"area,omitempty"`
	AvgDeposit            float64       `json:"avg_deposit,omitempty"`
	AvgPrice              float64       `json:"avg_price,omitempty"`
}

// RoomExtended - расширенная структура Room с добавление полей структуры Lot
type RoomExtended struct {
	ID                     int           `json:"id,omitempty"`
	LivingPlaces           []LivingPlace `json:"living_place,omitempty"`
	FlatID                 int           `json:"flat_id,omitempty"`
	MaxResidents           int           `json:"max_residents"`
	CurrNumberOfResidents  int           `json:"curr_number_of_residents"`
	Description            string        `json:"description,omitempty"`
	NumOfWindows           int           `json:"num_of_windows"`
	Balcony                bool          `json:"balcony"`
	NumOfTables            int           `json:"num_of_tables"`
	NumOfChairs            int           `json:"num_of_chairs"`
	TV                     bool          `json:"tv"`
	NumOFCupboards         int           `json:"num_of_cupboards"`
	Area                   int           `json:"area"`
	AvgDeposit             float64       `json:"avg_deposit,omitempty"`
	AvgPrice               float64       `json:"avg_price,omitempty"`
	Address                string        `json:"address,omitempty"`
	Floor                  int           `json:"floor,omitempty"`
	FloorsTotal            int           `json:"floor_total,omitempty"`
	MetroStation           string        `json:"metro,omitempty"`
	TimeToMetroByTransport int           `json:"ttmetro_transport,omitempty"`
	FlatArea               int           `json:"flat_area,omitempty"`
	Long                   float64       `json:"long,omitempty"`
	Lat                    float64       `json:"lat,omitempty"`
}
