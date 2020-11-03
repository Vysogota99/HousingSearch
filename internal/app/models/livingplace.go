package models

// LivingPlace - спальное место в комнате
type LivingPlace struct {
	ID          int     `json:"id,omitempty"`
	RoomID      int     `json:"room_id,omitempty"`
	ResidentID  int     `json:"resident_id,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Description string  `binding:"required" json:"description,omitempty"`
	NumOFBerth  int     `binding:"required" json:"num_of_berth,omitempty"`
	Deposit     float64 `json:"deposit,omitempty"`
}

// MapLp ...
var MapLp = map[string]int8{
	"avg_price":   1,
	"avg_deposit": 1,
}
