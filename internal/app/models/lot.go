package models

// Lot - готовое объявление
type Lot struct {
	ID                      int     `json:"id,omitempty"`
	OwnerID                 int     `json:"owner_id,omitempty"`
	Rooms                   []Room  `json:"rooms,omitempty"`
	Address                 string  `json:"address,omitempty"`
	Coordinates             Point   `json:"coordinates,omitempty"`
	Description             string  `json:"description,omitempty"`
	TimeToMetroONFoot       int     `json:"ttmetro_food,omitempty"`
	TimeToMetroByTransport  int     `json:"ttmetro_transport,omitempty"`
	MetroStation            string  `json:"metro,omitempty"`
	Floor                   int     `json:"floor,omitempty"`
	FloorsTotal             int     `json:"floor_total,omitempty"`
	Area                    int     `json:"area,omitempty"`
	Repairs                 int     `json:"repair,omitempty"`
	PassElevator            bool    `json:"passenger_elevator"`
	ServiceElevator         bool    `json:"serve_elevator"`
	Kitchen                 bool    `json:"kitchen"`
	MicrowaveOven           bool    `json:"microwave_oven"`
	Bathroom                int     `json:"bathroom,omitempty"`
	Refrigerator            bool    `json:"refrigerator"`
	Dishwasher              bool    `json:"dishwasher"`
	Stove                   int     `json:"stove,omitempty"`
	VacuumCleaner           bool    `json:"vacuum_cleaner"`
	Dryer                   bool    `json:"dryer"`
	Internet                bool    `json:"internet"`
	Animals                 bool    `json:"animals"`
	Smoking                 bool    `json:"smoking"`
	Heating                 int     `json:"heating,omitempty"`
	IsVisible               bool    `json:"is_visible,omitempty"`
	TotalNumberOfResidents  int     `json:"total_num_of_residents,omitempty"`
	CurrNumberOfResidents   int     `json:"curr_num_of_residents,omitempty"`
	Price                   float64 `json:"flat_price,omitempty"`
	Deposit                 float64 `json:"flat_deposit,omitempty"`
	AvgPricePerResident     float64 `json:"avg_price_per_resident,omitempty"`
	AvgPriceDepositResident float64 `json:"avg_deposit_per_resident,omitempty"`
	CreatedAt               string  `json:"created_at,omitempty"`
	UpdatedAt               string  `json:"updated_at,omitempty"`
	IsConstructor           bool
}

// FlatConstructor - конструктор квартиры
type FlatConstructor struct {
	ID                     int    `json:"id"`
	OwnerID                int    `json:"owner_id"`
	Rooms                  []Room `binding:"required" json:"rooms"`
	Address                string `binding:"required" json:"address"`
	Coordinates            Point  `binding:"required" json:"coordinates"`
	TimeToMetroONFoot      int    `binding:"required" json:"ttmetro_food"`
	TimeToMetroByTransport int    `binding:"required" json:"ttmetro_transport"`
	MetroStation           string `binding:"required" json:"metro"`
	Floor                  int    `binding:"required" json:"floor"`
	FloorsTotal            int    `binding:"required" json:"floor_total"`
	Area                   int    `binding:"required" json:"area"`
	Repairs                int    `binding:"required" json:"repair"`
	PassElevator           bool   `binding:"required" json:"passenger_elevator"`
	ServiceElevator        bool   `binding:"required" json:"serve_elevator"`
	Kitchen                bool   `binding:"required" json:"kitchen"`
	MicrowaveOven          bool   `binding:"required" json:"microwave_oven"`
	Bathroom               int    `binding:"required" json:"bathroom"`
	Refrigerator           bool   `binding:"required" json:"refrigerator"`
	Dishwasher             bool   `binding:"required" json:"dishwasher"`
	Stove                  int    `binding:"required" json:"stove"`
	VacuumCleaner          bool   `binding:"required" json:"vacuum_cleaner"`
	Dryer                  bool   `binding:"required" json:"dryer"`
	Internet               bool   `binding:"required" json:"internet"`
	Animals                bool   `binding:"required" json:"animals"`
	Smoking                bool   `binding:"required" json:"smoking"`
	TotalNumberOfResidents int    `binding:"required" json:"total_num_of_residents"`
	Heating                int    `binding:"required" json:"heating"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

// Point - ...
type Point struct {
	X      float64 `binding:"required" json:"lat"`
	Y      float64 `binding:"required" json:"lon"`
	CellID uint64
}

// Repair - тип ремонта в квартире
type Repair struct {
	ID          int
	Name        string
	Description string
}

// Stove - тип плиты
type Stove struct {
	ID          int
	Name        string
	Description string
}

// Heating - отопление
type Heating struct {
	ID          int
	Name        string
	Description string
}

// TestLot ...
var TestLot = Lot{
	OwnerID: 1,
	Rooms: []Room{
		Room{
			LivingPlaces: []LivingPlace{
				LivingPlace{
					ResidentID: 2,
					NumOFBerth: 1,
				},
				LivingPlace{
					ResidentID: 2,
					NumOFBerth: 1,
				},
			},
			MaxResidents:          2,
			CurrNumberOfResidents: 0,
			Windows:               true,
			Balcony:               false,
			NumOfTables:           2,
			NumOfChairs:           4,
			TV:                    false,
			Furniture:             true,
			Area:                  25,
		},
		Room{
			LivingPlaces: []LivingPlace{
				LivingPlace{
					ResidentID:  2,
					Price:       12000,
					Description: "В то время некий безымянный печатник создал большую коллекцию размеров и форм шрифтов, используя Lorem Ipsum для распечатки образцов.",
					NumOFBerth:  2,
					Deposit:     12000,
				},
			},
			MaxResidents:          2,
			CurrNumberOfResidents: 0,
			Description:           "Комната 2",
			Windows:               true,
			Balcony:               true,
			NumOfTables:           1,
			NumOfChairs:           1,
			TV:                    true,
			Furniture:             true,
			Area:                  20,
		},
	},
	Address: "Россия, Москва, Коломенский проезд, 23к1, кв32",
	Coordinates: Point{
		X: 55.667959,
		Y: 37.656157,
	},
	TimeToMetroONFoot:      15,
	TimeToMetroByTransport: 10,
	MetroStation:           "Коломенская",
	Floor:                  8,
	FloorsTotal:            12,
	Area:                   60,
	Repairs:                1,
	PassElevator:           true,
	ServiceElevator:        true,
	Kitchen:                true,
	MicrowaveOven:          true,
	Bathroom:               1,
	Refrigerator:           true,
	Dishwasher:             false,
	Stove:                  1,
	Dryer:                  true,
	VacuumCleaner:          false,
	Internet:               true,
	Animals:                false,
	Heating:                1,
	Smoking:                true,
}

// MapFlat - ...
var MapFlat = map[string]int8{
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
	"price":                  1,
	"deposit":                1,
}
