package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"userportal/internal/app/dto"

	"github.com/jmoiron/sqlx"
)

const insertUserQuery = `INSERT INTO users (first_name, last_name, email) VALUES (:first_name, :last_name, :email)`
const getAllUsersQuery = `SELECT * FROM users`
const getUserByEmailQuery = `SELECT * FROM users WHERE email = $1`
const updateUserQuery = `UPDATE users SET first_name = :first_name, last_name = :last_name WHERE email = :email`
const deleteUserQuery = `DELETE FROM users WHERE email = $1`

type ErrUserNotFound struct {
	Email string
}

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

// Implement the Error() method for the custom error type.
func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with email '%s' not found", e.Email)
}

func NewUserRepository(conn *sqlx.DB) UserRepository {
	return &userRepository{
		db: conn,
	}
}

func (ur userRepository) GetUserByEmail(email string) (*dto.User, error) {
	user := dto.User{}
	err := ur.db.Get(&user, getUserByEmailQuery, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrUserNotFound{Email: email}
		}
		return nil, err
	}
	return &user, nil

}
func (ur userRepository) CreateUser(user dto.User) error {
	isUserExisting, _ := isUserExisting(ur, user)
	if isUserExisting {
		return errors.New("Duplicate request , user already existing with email " + user.Email)
	}
	_, err := ur.db.NamedExec(insertUserQuery, user)

	return err
}

func isUserExisting(ur userRepository, user dto.User) (bool, error) {
	existingUser, err1 := ur.GetUserByEmail(user.Email)
	_, isUserNotFoundErr := err1.(*ErrUserNotFound)
	if err1 != nil && !isUserNotFoundErr {

		log.Fatal("Error while retreiving user", err1)
	}
	if existingUser != nil {
		return true, nil
	}
	return false, nil
}

func (ur userRepository) CreateUsers(users []dto.User) error {
	tx, err := ur.db.Beginx()
	if err != nil {
		return err
	}

	for _, user := range users {
		isUserExisting, _ := isUserExisting(ur, user)
		if isUserExisting {
			return errors.New("Duplicate request , user already existing with email " + user.Email)
		}

		if _, err := tx.NamedExec(insertUserQuery, user); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (ur userRepository) GetAllUsers() ([]dto.User, error) {
	users := []dto.User{}
	err := ur.db.Select(&users, getAllUsersQuery)
	if err != nil {
		fmt.Println("Error in db ", err)
		return nil, err
	}
	return users, nil
}

func (ur userRepository) UpdateUser(user dto.User) error {
	isUserExisting, _ := isUserExisting(ur, user)

	if !isUserExisting {
		return &ErrUserNotFound{Email: user.Email}
	}

	_, err := ur.db.NamedExec(updateUserQuery, user)
	return err

}

func (ur userRepository) DeleteUserByEmail(email string) error {

	result, err := ur.db.Exec(deleteUserQuery, email)

	if rowsEffected, _ := result.RowsAffected(); err == nil && rowsEffected == 0 {
		fmt.Println("rows effected", rowsEffected)
		return &ErrUserNotFound{Email: email}
	}
	return err
}
