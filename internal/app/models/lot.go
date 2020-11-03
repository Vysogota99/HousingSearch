package models

// Lot - жилье для сдачи
type Lot struct {
	ID                      int     `json:"id,omitempty"`
	OwnerID                 int     `json:"owner_id,omitempty"`
	Rooms                   []Room  `binding:"required" json:"rooms,omitempty"`
	Address                 string  `binding:"required" json:"address,omitempty"`
	Coordinates             Point   `binding:"required" json:"coordinates,omitempty"`
	Description             string  `binding:"required" json:"description,omitempty"`
	TimeToMetroONFoot       int     `binding:"required" json:"ttmetro_food,omitempty"`
	TimeToMetroByTransport  int     `binding:"required" json:"ttmetro_transport,omitempty"`
	MetroStation            string  `binding:"required" json:"metro,omitempty"`
	Floor                   int     `binding:"required" json:"floor,omitempty"`
	FloorsTotal             int     `binding:"required" json:"floor_total,omitempty"`
	Area                    int     `binding:"required" json:"area,omitempty"`
	Repairs                 int     `binding:"required" json:"repair,omitempty"`
	Elevators               bool    `binding:"required" json:"elevator,omitempty"`
	Bathroom                int     `binding:"required" json:"bathroom,omitempty"`
	Refrigerator            bool    `binding:"required" json:"refrigerator,omitempty"`
	Dishwasher              bool    `binding:"required" json:"dishwasher,omitempty"`
	GasStove                bool    `binding:"required" json:"gasStove,omitempty"`
	ElectricStove           bool    `binding:"required" json:"electric_stove,omitempty"`
	VacuumCleaner           bool    `json:"vacuumCleaner,omitempty"`
	Internet                bool    `json:"internet,omitempty"`
	Animals                 bool    `json:"animals,omitempty"`
	Smoking                 bool    `json:"smoking,omitempty"`
	IsVisible               bool    `json:"is_visible,omitempty"`
	TotalNumberOfResidents  int     `json:"total_num_of_residents,omitempty"`
	CurrNumberOfResidents   int     `json:"curr_num_of_residents,omitempty"`
	Price                   float64 `json:"flat_price,omitempty"`
	AvgPricePerResident     float64 `json:"avg_price_per_resident,omitempty"`
	AvgPriceDepositResident float64 `json:"avg_deposit_per_resident,omitempty"`
	CreatedAt               string  `json:"created_at,omitempty"`
	Deposit                 float64 `json:"flat_deposit,omitempty"`
}

// Point - ...
type Point struct {
	X float64 `binding:"required" json:"lat"`
	Y float64 `binding:"required" json:"lon"`
}

// TestLot ...
var TestLot = Lot{
	OwnerID: 1,
	Rooms: []Room{
		Room{
			LivingPlaces: []LivingPlace{
				LivingPlace{
					ResidentID:  2,
					Price:       5000,
					Description: "В то время некий безымянный печатник создал большую коллекцию размеров и форм шрифтов, используя Lorem Ipsum для распечатки образцов.",
					NumOFBerth:  1,
					Deposit:     5000,
				},
				LivingPlace{
					ResidentID:  2,
					Price:       5000,
					Description: "Здесь ваш текст.. Многие программы электронной вёрстки и редакторы HTML используют Lorem Ipsum в качестве текста по умолчанию.",
					NumOFBerth:  1,
					Deposit:     5000,
				},
			},
			MaxResidents:          2,
			CurrNumberOfResidents: 0,
			Description:           "Комната 1",
			NumOfWindows:          1,
			Balcony:               false,
			NumOfTables:           2,
			NumOfChairs:           4,
			TV:                    false,
			NumOFCupboards:        1,
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
			NumOfWindows:          1,
			Balcony:               true,
			NumOfTables:           1,
			NumOfChairs:           1,
			TV:                    true,
			NumOFCupboards:        1,
			Area:                  20,
		},
	},
	Address: "Россия, Москва, Коломенский проезд, 23к1, кв32",
	Coordinates: Point{
		X: 55.667959,
		Y: 37.656157,
	},
	Description:            "Впервые сдаётся светлая 1-комнатная квартира с евроремонтом, смежным санузлом и балконом в ЖК Одинбург, без депозита! Полы  ламинат и плитка. Квартира оборудована кухонным гарнитуром, мебелью и бытовой техникой. Есть кондиционер, холодильник, телевизор, плита, душ, посудомоечная и стиральная машины. Проведён интернет. Окна выходят на улицу, есть бесплатная парковка.Собственник готов заселить не более 2 взрослых жильцов или заключить договор с юридическим лицом. Гражданство: РФ. С домашним питомцем нельзя. Дом находится в районе с развитой инфраструктурой: в пешей доступности  школы и детские сады, продуктовые супермаркеты, аптеки и различные медицинские учреждения, салоны красоты и парикмахерские, фитнес-клубы, кафе и закусочные. До МЦД-1 Одинцово  14 минут на общественном транспорте, до м. Кунцевская  40 минут. Автомобилистам будет удобно выезжать на МКАД через Минское шоссе",
	TimeToMetroONFoot:      15,
	TimeToMetroByTransport: 10,
	MetroStation:           "Коломенская",
	Floor:                  8,
	FloorsTotal:            12,
	Area:                   60,
	Repairs:                1,
	Elevators:              true,
	Bathroom:               1,
	Refrigerator:           true,
	Dishwasher:             false,
	GasStove:               true,
	ElectricStove:          false,
	VacuumCleaner:          false,
	Internet:               true,
	Animals:                false,
	Smoking:                true,
	IsVisible:              true,
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
