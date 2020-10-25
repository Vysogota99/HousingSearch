package models

// LivingPlace - спальное место в комнате
type LivingPlace struct {
	ID          int     `json:"id,omitempty"`
	RoomID      int     `json:"room_id,omitempty"`
	ResidentID  int     `json:"resident_id,omitempty"`
	Price       float64 `binding:"required" json:"price,omitempty"`
	Description string  `binding:"required" json:"description,omitempty"`
	NumOFBerth  int     `binding:"required" json:"num_of_berth,omitempty"`
	Deposit     float64 `binding:"required" json:"deposit,omitempty"`
}
