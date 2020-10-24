package models

// LivingPlace - спальное место в комнате
type LivingPlace struct {
	ID          int
	RoomID      int
	ResidentID  int
	Price       float64
	Description string
	NumOFBerth  int
	Deposit     float64
}
