package models

// Roles - роли пользователей в системе
type Roles struct {
	ID          string `JSON:"id"`
	Name        string `JSON:"name"`
	Description string `JSON:"description"`
}
