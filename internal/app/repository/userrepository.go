package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"userportal/internal/app/dto"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user dto.User) error
	CreateUsers(users []dto.User) error
	GetAllUsers() ([]dto.User, error)
	GetUserByEmail(email string) (*dto.User, error)
	UpdateUser(user dto.User) error
	DeleteUserByEmail(email string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) UserRepository {
	return &userRepository{
		db: conn,
	}
}

func (ur userRepository) GetUserByEmail(email string) (*dto.User, error) {
	user := dto.User{}
	err := ur.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil

}
func (ur userRepository) CreateUser(user dto.User) error {
	isUserExisting, err := isUserExisting(ur, user)
	fmt.Println("is user existing", isUserExisting)
	if isUserExisting {
		return errors.New("Duplicate request , user already existing with email" + user.Email)
	}
	if err != nil {
		return err
	}

	query := `INSERT INTO users (first_name, last_name, email) VALUES (:first_name,:last_name,:email)`

	_, err = ur.db.NamedExec(query, user)

	return err
}

func isUserExisting(ur userRepository, user dto.User) (bool, error) {
	existingUser, err1 := ur.GetUserByEmail(user.Email)
	fmt.Println("existung users", existingUser)
	if err1 != nil && err1.Error() != "user not found" {
		return false, err1
	}
	if existingUser != nil {
		fmt.Println("user already existing")
		return true, nil
	}
	return false, nil
}

func (ur userRepository) CreateUsers(users []dto.User) error {
	tx, err := ur.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (first_name, last_name, email) VALUES (:first_name,:last_name,:email)`
	for _, user := range users {
		isUserExisting, err := isUserExisting(ur, user)
		if isUserExisting {
			return errors.New("Duplicate request , user already existing with email" + user.Email)
		}
		if err != nil {
			return err
		}
		if _, err := tx.NamedExec(query, &user); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (ur userRepository) GetAllUsers() ([]dto.User, error) {
	users := []dto.User{}
	err := ur.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		fmt.Println("Error in db ", err)
		return nil, err
	}
	return users, nil
}

func (ur userRepository) UpdateUser(user dto.User) error {
	isUserExisting, err := isUserExisting(ur, user)
	fmt.Println("is user existing", isUserExisting)
	if !isUserExisting {
		return errors.New("no user exists with given email")
	}
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET first_name = :first_name, last_name = :last_name
		WHERE email = :email
	`

	_, err = ur.db.NamedExec(query, user)
	return err

}

func (ur userRepository) DeleteUserByEmail(email string) error {
	query := `
		DELETE FROM users
		WHERE email = $1
	`

	result, err := ur.db.Exec(query, email)

	if rowsEffected, _ := result.RowsAffected(); err == nil && rowsEffected == 0 {
		fmt.Println("rows effected", rowsEffected)
		return errors.New("no user exists with given email")
	}
	return err
}
