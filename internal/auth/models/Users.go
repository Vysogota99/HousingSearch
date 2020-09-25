package models

// User ...
type User struct {
	ID                 int    `JSON:"id"`
	TelephoneNumber    string `JSON:"telephone_number"`
	TelegramProfile    string `JSON:"telegram_profile"`
	Email              string `JSON:"email"`
	Role               int    `JSON:"role"`
	PassSeries         string `JSON:"passport_series"`
	PassNumber         string `JSON:"passport_number"`
	PassDateOfIssue    string `JSON:"passport_date_of_issue"`
	PassDepartmentCode string `JSON:"passport_department_code"`
	PassIssueBy        string `JSON:"passport_issue_by"`
	PassName           string `JSON:"passport_name"`
	PassLastName       string `JSON:"passport_lastname"`
	PassPatronymic     string `JSON:"passport_patronic"`
	PassSex            string `JSON:"passport_sex"`
	PassDateOfBirth    string `JSON:"passport_date_of_birth"`
	PassPlaceOfBirth   string `JSON:"passport_place_of_birth"`
	PassRegistration   string `JSON:"passport_registration"`
}

// Roles - роли пользователей в системе
type Roles struct {
	ID          string `JSON:"id"`
	Name        string `JSON:"name"`
	Description string `JSON:"description"`
}
