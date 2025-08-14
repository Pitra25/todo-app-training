package service

import (
	"todo-app/internal/repository"
	"todo-app/internal/repository/mysql/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	GenerateRefrachToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Users interface {
	GetUserById(userId int) (*models.UserResponse, error)
	GetUserAll() (*[]models.UserResponse, error)
	UpdateUser(userId int, input *models.UpdateUserInpur) error
	DeleteUser(userId int) error
}

type TodoList interface {
	Create(userId int, list models.TodoList) (int, error)
	GetAll(userId int) ([]models.TodoList, error)
	GetById(userId, listId int) (models.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input models.UpdadeListInput) error
}

type TodoItems interface {
	Create(userId, listId int, item models.TodoItems) (int, error)
	GetAllItemsList(userId, listId int) ([]models.TodoItems, error)
	GetAllItem() ([]models.TodoItems, error)
	GetById(userId, itemId int) (models.TodoItems, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input models.UpdadeItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItems
	Users
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItems:     NewTodoItemsService(repos.TodoItems, repos.TodoList),
		Users:         NewUserService(repos.Users),
	}
}
