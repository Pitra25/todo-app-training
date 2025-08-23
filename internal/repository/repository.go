package repository

import (
	"todo-app/internal/repository/mysql"
	"todo-app/internal/repository/mysql/methods"
	"todo-app/internal/repository/mysql/models"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
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
	Create(listId int, item models.TodoItems) (int, error)
	GetAllItemsList(userId, listId int) ([]models.TodoItems, error)
	GetAllItem() ([]models.TodoItems, error)
	GetById(userId, listId int) (models.TodoItems, error)
	Delete(userId, itemId int) error
	Update(userId, listId int, input models.UpdateItemInput) error
}

type MySql interface {
	NewMySqlDB(cfg *mysql.ConfigMySql) (*sqlx.DB, error)
}

type Emails interface {
	SaveCodeUser(code string, userId int) error
	GetCodeUser(userId int) (methods.ResponseCode, error)
	UpdateStatusUser(userId int) error
	DeleteRecord(id, userId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItems
	Users
	Emails
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Authorization: methods.NewAuthMySql(db),
		TodoList:      methods.NewTodoListMySql(db, rdb),
		TodoItems:     methods.NewTodoItemsMySql(db, rdb),
		Users:         methods.NewUsersMySQL(db),
		Emails:        methods.NewEmailMySql(db, rdb),
	}
}
