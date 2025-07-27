package repository

import (
	"fmt"
	"strings"
	"todo-app/types"

	"github.com/jmoiron/sqlx"
)

type TodoListMySql struct {
	db *sqlx.DB
}

func NewTodoListMySql(db *sqlx.DB) *TodoListMySql {
	return &TodoListMySql{db: db}
}

func (r *TodoListMySql) Create(userId int, list types.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", todoListsTable)
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

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES (?, ?)", usersListsTable)
	_, err = tx.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), tx.Commit()
}

func (r *TodoListMySql) GetAll(userId int) ([]types.TodoList, error) {
	var lists []types.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = ?",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListMySql) GetById(userId, listId int) (types.TodoList, error) {
	var lists types.TodoList

	storageRecords, err := Get(listId, "list")
	if err != nil {
		return types.TodoList{}, err
	}
	if storageRecords.list.Title != "" {
		return storageRecords.list, nil
	}

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?",
		todoListsTable, usersListsTable)
	err = r.db.Get(&lists, query, userId, listId)
	if err != nil {
		return lists, err
	}

	err = Create(RecordingStruct{
		ID:   lists.Id,
		list: lists,
	})

	return lists, err
}

func (r *TodoListMySql) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE tl FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListMySql) Update(userId, listId int, input types.UpdadeListInput) error {
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

	query := fmt.Sprintf("UPDATE %s tl JOIN %s ul ON tl.id = ul.list_id SET %s WHERE ul.list_id = ? AND ul.user_id = ?",
		todoListsTable, usersListsTable, setQuery)
	args = append(args, listId, userId)

	_, err := r.db.Exec(query, args...)
	return err
}
