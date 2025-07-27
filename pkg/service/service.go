package service

import (
	"todo-app/pkg/repository"
	"todo-app/types"
)

type Authorization interface {
	CreateUser(user types.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list types.TodoList) (int, error)
	GetAll(userId int) ([]types.TodoList, error)
	GetById(userId, listId int) (types.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input types.UpdadeListInput) error
}

type TodoItems interface {
	Create(userId, listId int, item types.TodoItems) (int, error)
	GetAllItemsList(userId, listId int) ([]types.TodoItems, error)
	GetAllItem() ([]types.TodoItems, error)
	GetById(userId, itemId int) (types.TodoItems, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input types.UpdadeItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItems
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItems:     NewTodoItemsService(repos.TodoItems, repos.TodoList),
	}
}
