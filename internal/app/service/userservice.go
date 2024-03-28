package service

import (
	"userportal/internal/app/dto"
	"userportal/internal/app/repository"
)

type UserService interface {
	CreateUser(user dto.User) error
	CreateUsers(users []dto.User) error
	GetAllUsers() ([]dto.User, error)
	GetUserByEmail(email string) (*dto.User, error)
	UpdateUser(user dto.User) error
	DeleteUserByEmail(email string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func (us *userService) CreateUser(user dto.User) error {
	return us.userRepo.CreateUser(user)

}

func (us *userService) CreateUsers(users []dto.User) error {
	return us.userRepo.CreateUsers(users)
}

func (us *userService) GetAllUsers() ([]dto.User, error) {
	return us.userRepo.GetAllUsers()
}

// GetUserByEmail implements UserService.
func (us *userService) GetUserByEmail(email string) (*dto.User, error) {
	return us.userRepo.GetUserByEmail(email)
}

func (us *userService) UpdateUser(user dto.User) error {
	return us.userRepo.UpdateUser(user)

}
func (us *userService) DeleteUserByEmail(email string) error {
	return us.userRepo.DeleteUserByEmail(email)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}
