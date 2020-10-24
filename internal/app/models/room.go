package models

// Room - структура описывающая комнату в квартире
type Room struct {
	ID                    int
	LivingPlaces          []LivingPlace
	FlatID                int
	MaxResidents          int
	CurrNumberOfResidents int
	Description           string
	NumOfWindows          int
	Balcony               bool
	NumOfTables           int
	NumOfChairs           int
	TV                    bool
	NumOFCupboards        int
	Area                  int
	AvgDeposit            float64
	AvgPrice              float64
}
