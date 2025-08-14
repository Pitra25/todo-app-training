package service

import (
	"todo-app/internal/repository"
	"todo-app/internal/repository/mysql/models"
)

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(userId int) (*models.UserResponse, error) {
	return s.repo.GetUserById(userId)
}

func (s *UserService) GetUserAll() (*[]models.UserResponse, error) {
	return s.repo.GetUserAll()
}

func (s *UserService) UpdateUser(userId int, input *models.UpdateUserInpur) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateUser(userId, input)
}

func (s *UserService) DeleteUser(userId int) error {
	return s.repo.DeleteUser(userId)
}
