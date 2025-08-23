package service

import (
	"todo-app/internal/repository"
	"todo-app/internal/repository/mysql/models"
	"todo-app/internal/service/methods"
	"todo-app/pkg/email"

	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	GenerateRefreshToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Users interface {
	GetUserById(userId int) (*models.UserResponse, error)
	GetUserAll() (*[]models.UserResponse, error)
	UpdateUser(userId int, input *models.UpdateUserInput) error
	DeleteUser(userId int) error
}

type TodoList interface {
	Create(userId int, list models.TodoList) (int, error)
	GetAll(userId int) ([]models.TodoList, error)
	GetById(userId, listId int) (models.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input models.UpdateListInput) error
}

type TodoItems interface {
	Create(userId, listId int, item models.TodoItems) (int, error)
	GetAllItemsList(userId, listId int) ([]models.TodoItems, error)
	GetAllItem() ([]models.TodoItems, error)
	GetById(userId, itemId int) (models.TodoItems, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input models.UpdateItemInput) error
}

type Emails interface {
	SendEmail(to string, userId int) error
	ConfirmationEmail(code string, userId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItems
	Users
	Emails
}

func NewService(
	repos *repository.Repository,
	eCient *email.Email,
	rdb *redis.Client,
) *Service {
	return &Service{
		Authorization: methods.NewAuthService(repos.Authorization),
		TodoList:      methods.NewTodoListService(repos.TodoList),
		TodoItems:     methods.NewTodoItemsService(repos.TodoItems, repos.TodoList),
		Users:         methods.NewUserService(repos.Users),
		Emails:        methods.NewEmailService(repos.Emails, eCient, rdb),
	}
}
