package models

// Room - структура описывающая комнату в квартире
type Room struct {
	ID                    int           `json:"id,omitempty"`
	LivingPlaces          []LivingPlace `binding:"required" json:"living_place,omitempty"`
	FlatID                int           `json:"flat_id,omitempty"`
	MaxResidents          int           `binding:"required" json:"max_residents,omitempty"`
	CurrNumberOfResidents int           `json:"curr_number_of_residents,omitempty,omitempty"`
	Description           string        `binding:"required" json:"description,omitempty"`
	NumOfWindows          int           `binding:"required" json:"num_of_windows,omitempty"`
	Balcony               bool          `binding:"required" json:"balcony,omitempty"`
	NumOfTables           int           `binding:"required" json:"num_of_tables,omitempty"`
	NumOfChairs           int           `binding:"required" json:"num_of_chairs,omitempty"`
	TV                    bool          `binding:"required" json:"tv,omitempty"`
	NumOFCupboards        int           `binding:"required" json:"num_of_cupboards,omitempty"`
	Area                  int           `binding:"required" json:"area,omitempty"`
	AvgDeposit            float64       `json:"avg_deposit,omitempty,omitempty"`
	AvgPrice              float64       `json:"avg_price,omitempty,omitempty"`
}
