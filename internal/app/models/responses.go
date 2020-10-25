package models

// Paginations ...
type Paginations struct {
	Data        []Lot `json:"data"`
	CurrentPage int   `json:"curr_page"`
	NumPages    int   `json:"num_pages"`
}
