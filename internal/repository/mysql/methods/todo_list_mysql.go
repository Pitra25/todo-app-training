package methods

import (
	"fmt"
	"strings"
	"todo-app/internal/repository/mysql/models"
	storage "todo-app/pkg/cache/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TodoListMySql struct {
	db  *sqlx.DB
	rdb *storage.RedisCLientDB
}

func NewTodoListMySql(db *sqlx.DB, rdb *redis.Client) *TodoListMySql {
	return &TodoListMySql{
		db:  db,
		rdb: &storage.RedisCLientDB{Db: rdb},
	}
}

func (r *TodoListMySql) Create(userId int, list models.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES (?, ?)",
		models.TodoListsTable,
	)
	row, err := tx.Exec(createListQuery, list.Title, list.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) VALUES (?, ?)",
		models.UsersListsTable,
	)
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), tx.Commit()
}

func (r *TodoListMySql) GetAll(userId int) ([]models.TodoList, error) {
	var lists []models.TodoList

	query := fmt.Sprintf(
		"SELECT id, title, description FROM %s",
		models.TodoListsTable)
	err := r.db.Select(&lists, query)

	return lists, err
}

func (r *TodoListMySql) GetById(userId, listId int) (models.TodoList, error) {
	var lists models.TodoList

	// Check Redis cache first
	storageRecords, err := r.rdb.Get(listId, storage.List)
	if err != nil {
		logrus.Error("method: GetById list.", err.Error())
	} else if storageRecords.List.Title != "" && storageRecords != nil {
		logrus.Debug("GetById list from redis cache. list id:", listId)
		return storageRecords.List, nil
	}

	// Query the database for the list
	query := fmt.Sprintf(
		"SELECT tl.id, tl.title, tl.description FROM %s tl "+
			"INNER JOIN %s ul on tl.id = ul.list_id "+
			"WHERE ul.user_id = ? AND ul.list_id = ?",
		models.TodoListsTable, models.UsersListsTable,
	)

	err = r.db.Get(&lists, query, userId, listId)
	if err != nil {
		return lists, err
	}

	// Save to Redis cache
	if err := r.rdb.Create(&storage.Recording{
		ID:   lists.Id,
		List: lists,
	}); err != nil {
		logrus.Error("err create record redis.", err.Error())
		return lists, nil
	}

	return lists, nil
}

func (r *TodoListMySql) Delete(userId, listId int) error {
	query := fmt.Sprintf(
		"DELETE tl FROM %s tl"+
			"INNER JOIN %s ul ON tl.id = ul.list_id"+
			"WHERE ul.user_id = ? AND ul.list_id = ?",
		models.TodoListsTable, models.UsersListsTable,
	)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListMySql) Update(userId, listId int, input models.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Title != nil {
		setValues = append(setValues, "tl.title = ?")
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		setValues = append(setValues, "tl.description = ?")
		args = append(args, *input.Description)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		"UPDATE %s tl JOIN %s ul ON tl.id = ul.list_id SET %s"+
			"WHERE ul.list_id = ? AND ul.user_id = ?",
		models.TodoListsTable, models.UsersListsTable,
		setQuery,
	)
	args = append(args, listId, userId)

	_, err := r.db.Exec(query, args...)
	return err
}
