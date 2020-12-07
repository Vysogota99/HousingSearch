package postgresstore

import (
	"database/sql"

	"context"

	"fmt"

	"github.com/Vysogota99/HousingSearch/pkg/authService"
	_ "github.com/lib/pq"
)

// UserRepository - реализует функционал модели User
type UserRepository struct {
	store *StorePSQL
}

// User model
type User struct {
	ID                 int64
	TelephoneNumber    string
	TelegramProfile    sql.NullString
	VkProfile          sql.NullString
	Email              sql.NullString
	Role               int32
	PassSeries         sql.NullString
	PassNumber         sql.NullString
	PassDateOfIssue    sql.NullString
	PassDepartmentCode sql.NullString
	PassIssueBy        sql.NullString
	PassName           string
	PassLastName       string
	PassPatronymic     sql.NullString
	PassSex            string
	PassDateOfBirth    string
	PassPlaceOfBirth   sql.NullString
	PassRegistration   sql.NullString
	Password           string
	AvatarPath         string
}

// CreateUser - создает нового пользователя в базе данных
func (u *UserRepository) CreateUser(user *authService.User) error {
	db, err := sql.Open("postgres", u.store.ConnString)
	if err != nil {
		return err
	}
	defer db.Close()

	pass := HashPassword(user.Password)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, "INSERT INTO users(telephone_number, role, password, avatar_path) VALUES ($1, $2, $3, $4) RETURNING id;", user.TelephoneNumber, user.Role, <-pass, user.AvatartPath).Scan(&user.ID)
	if row != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %w", rollbackErr)
		}
		return row
	}

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %w", rollbackErr)
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO passports_info(user_id, passp_name, passp_lastname, passp_sex, passp_date_of_birth) VALUES ($1, $2, $3, $4, $5);", user.ID, user.PassName, user.PassLastName, user.PassSex, user.PassDateOfBirth)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %w", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetUser - получение пользователя (по номеру телефона) из базы данных
func (u *UserRepository) GetUser(telephoneNumber string) (*authService.User, error) {
	db, err := sql.Open("postgres", u.store.ConnString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	user := &User{}
	row := db.QueryRow(`
							SELECT u.id, u.email, u.vk_profile, u.telephone_number, u.role, u.password, u.avatar_path, p.passp_series, p.passp_number,
							p.passp_date_of_issue, p.passp_department_code, p.passp_issue_by, p.passp_name,
							p.passp_lastname, p.passp_patronymic, p.passp_sex, p.passp_date_of_birth,
							p.passp_place_of_birth, p.passp_registration
							FROM users AS u
							JOIN passports_info AS p
							ON u.id = p.user_id
							WHERE u.telephone_number = $1;
						`, telephoneNumber).Scan(
		&user.ID,
		&user.Email,
		&user.VkProfile,
		&user.TelephoneNumber,
		&user.Role,
		&user.Password,
		&user.AvatarPath,
		&user.PassSeries,
		&user.PassNumber,
		&user.PassDateOfIssue,
		&user.PassDepartmentCode,
		&user.PassIssueBy,
		&user.PassName,
		&user.PassLastName,
		&user.PassPatronymic,
		&user.PassSex,
		&user.PassDateOfBirth,
		&user.PassPlaceOfBirth,
		&user.PassRegistration,
	)

	switch {
	case row == sql.ErrNoRows:
		return nil, fmt.Errorf("No user with telephone number %s", telephoneNumber)
	case row != nil:
		return nil, row
	default:
		resultUser := &authService.User{}
		userCopy(user, resultUser)
		return resultUser, nil
	}
}

func userCopy(from *User, to *authService.User) {
	to.ID = from.ID
	to.TelephoneNumber = from.TelephoneNumber
	to.Role = from.Role
	to.PassName = from.PassName
	to.PassLastName = from.PassLastName
	to.PassSex = from.PassSex
	to.PassDateOfBirth = from.PassDateOfBirth
	to.Password = from.Password
	to.AvatartPath = from.AvatarPath

	if from.TelegramProfile.Valid {
		to.TelegramProfile = from.TelegramProfile.String
	}

	if from.VkProfile.Valid {
		to.VkProfile = from.VkProfile.String
	}

	if from.Email.Valid {
		to.Email = from.Email.String
	}

	if from.PassSeries.Valid {
		to.PassSeries = from.PassSeries.String
	}

	if from.PassNumber.Valid {
		to.PassNumber = from.PassNumber.String
	}

	if from.PassDateOfIssue.Valid {
		to.PassDateOfIssue = from.PassDateOfIssue.String
	}

	if from.PassDepartmentCode.Valid {
		to.PassDepartmentCode = from.PassDepartmentCode.String
	}

	if from.PassIssueBy.Valid {
		to.PassIssueBy = from.PassIssueBy.String
	}

	if from.PassPatronymic.Valid {
		to.PassPatronymic = from.PassPatronymic.String
	}

	if from.PassPlaceOfBirth.Valid {
		to.PassPlaceOfBirth = from.PassPlaceOfBirth.String
	}

	if from.PassRegistration.Valid {
		to.PassRegistration = from.PassRegistration.String
	}
}
