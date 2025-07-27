package repository

import (
	"todo-app/types"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user types.User) (int, error)
	GetUser(username, password string) (types.User, error)
}

type TodoList interface {
	Create(userId int, list types.TodoList) (int, error)
	GetAll(userId int) ([]types.TodoList, error)
	GetById(userId, listId int) (types.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input types.UpdadeListInput) error
}

type TodoItems interface {
	Create(listId int, item types.TodoItems) (int, error)
	GetAllItemsList(userId, listId int) ([]types.TodoItems, error)
	GetAllItem() ([]types.TodoItems, error)
	GetById(userId, listId int) (types.TodoItems, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input types.UpdadeItemInput) error
}

type Redis interface {
	Create(record RecordingStruct) error
	Get(key string) (RecordingStruct, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItems
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMySql(db),
		TodoList:      NewTodoListMySql(db),
		TodoItems:     NewTodoItemsMySql(db),
	}
}
