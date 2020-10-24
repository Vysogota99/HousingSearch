package models

// Lot - жилье для сдачи
type Lot struct {
	ID                      int     `JSON:""`
	OwnerID                 int     `JSON:""`
	Rooms                   []Room  `JSON:""`
	Address                 string  `JSON:""`
	Coordinates             Point   `JSON:""`
	Description             string  `JSON:""`
	Deposit                 float64 `JSON:""`
	TimeToMetroONFoot       int     `JSON:""`
	TimeToMetroByTransport  int     `JSON:""`
	MetroStation            string  `JSON:""`
	Floor                   int     `JSON:""`
	FloorsTotal             int     `JSON:""`
	Area                    int     `JSON:""`
	Repairs                 int     `JSON:""`
	Elevators               bool    `JSON:""`
	Bathroom                int     `JSON:""`
	Refrigerator            bool    `JSON:""`
	Dishwasher              bool    `JSON:""`
	GasStove                bool    `JSON:""`
	ElectricStove           bool    `JSON:""`
	VacuumCleaner           bool    `JSON:""`
	Internet                bool    `JSON:""`
	Animals                 bool    `JSON:""`
	Smoking                 bool    `JSON:""`
	IsVisible               bool    `JSON:""`
	TotalNumberOfResidents  int     `JSON:""`
	CurrNumberOfResidents   int     `JSON:""`
	AvgPricePerResident     float64 `JSON:""`
	AvgPriceDepositResident float64 `JSON:""`
}

// Point - ...
type Point struct {
	X float64 `json:"lat"`
	Y float64 `json:"lon"`
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
		{
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
	Description: "Впервые сдаётся светлая 1-комнатная квартира с евроремонтом, смежным санузлом и балконом в ЖК Одинбург, без депозита! Полы  ламинат и плитка. Квартира оборудована кухонным гарнитуром, мебелью и бытовой техникой. Есть кондиционер, холодильник, телевизор, плита, душ, посудомоечная и стиральная машины. Проведён интернет. Окна выходят на улицу, есть бесплатная парковка.Собственник готов заселить не более 2 взрослых жильцов или заключить договор с юридическим лицом. Гражданство: РФ. С домашним питомцем нельзя. Дом находится в районе с развитой инфраструктурой: в пешей доступности  школы и детские сады, продуктовые супермаркеты, аптеки и различные медицинские учреждения, салоны красоты и парикмахерские, фитнес-клубы, кафе и закусочные. До МЦД-1 Одинцово  14 минут на общественном транспорте, до м. Кунцевская  40 минут. Автомобилистам будет удобно выезжать на МКАД через Минское шоссе",

	Deposit:                22000,
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
