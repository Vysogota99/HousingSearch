package models

// Paginations ...
type Paginations struct {
	Data        []Lot `json:"data"`
	CurrentPage int   `json:"curr_page"`
	NumPages    int   `json:"num_pages"`
}

// PaginationsRoom ...
type PaginationsRoom struct {
	Data        []RoomExtended `json:"data"`
	CurrentPage int            `json:"curr_page"`
	NumPages    int            `json:"num_pages"`
}
