package postgres

import (
	_ "github.com/lib/pq"
)

// StorePSQL - реализует взаимодействие с базой данных
type StorePSQL struct {
	ConnString string
}

// New - инициализирует Store
func New(connString string) *StorePSQL {
	return &StorePSQL{
		ConnString: connString,
	}
}
