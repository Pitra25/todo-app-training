package repository

import (
	"fmt"
	"strings"
	"todo-app/types"

	"github.com/jmoiron/sqlx"
)

type TodoItemsMySql struct {
	db *sqlx.DB
}

func NewTodoItemsMySql(db *sqlx.DB) *TodoItemsMySql {
	return &TodoItemsMySql{db: db}
}

func (r *TodoItemsMySql) Create(listId int, item types.TodoItems) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", TodoItemssTable)
	row, err := tx.Exec(createItemQuery, item.Title, item.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	itemid, err := row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemssQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES (?, ?)", listsItemsTable)
	_, err = tx.Exec(createListItemssQuery, listId, itemid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(itemid), tx.Commit()
}

func (r *TodoItemsMySql) GetAllItemsList(userId, listId int) ([]types.TodoItems, error) {
	var items []types.TodoItems

	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE li.list_id = ? AND ul.user_id = ?",
		TodoItemssTable, listsItemsTable, usersListsTable,
	)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemsMySql) GetAllItem() ([]types.TodoItems, error) {
	var items []types.TodoItems

	query := fmt.Sprintf(
		"SELECT id, title, description, done FROM %s",
		TodoItemssTable,
	)
	if err := r.db.Select(&items, query); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemsMySql) GetById(userId, itemId int) (types.TodoItems, error) {
	var item types.TodoItems

	storageRecords, err := Get(itemId, "list")
	if err != nil {
		return types.TodoItems{}, err
	}
	if storageRecords.list.Title != "" {
		return storageRecords.items, nil
	}

	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE ti.id = ? AND ul.user_id = ?",
		TodoItemssTable, listsItemsTable, usersListsTable,
	)
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	err = Create(RecordingStruct{
		ID:    item.Id,
		items: item,
	})

	return item, err
}

func (r *TodoItemsMySql) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		"DELETE ti FROM %s ti "+
			"INNER JOIN %s li ON ti.id = li.item_id "+
			"INNER JOIN %s ul ON li.list_id = ul.list_id "+
			"WHERE ul.user_id = ? AND ti.id = ?",
		TodoItemssTable,
		listsItemsTable,
		usersListsTable,
	)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *TodoItemsMySql) Update(userId, itemId int, input types.UpdadeItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Title != nil {
		setValues = append(setValues, "ti.title = ?")
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		setValues = append(setValues, "ti.description = ?")
		args = append(args, *input.Description)
	}

	if input.Done != nil {
		setValues = append(setValues, "ti.done = ?")
		args = append(args, *input.Done)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		"UPDATE %s ti "+
			"INNER JOIN %s li ON ti.id = li.item_id "+
			"INNER JOIN %s ul ON li.list_id = ul.list_id "+
			"SET %s "+
			"WHERE ul.user_id = ? AND ti.id = ?",
		TodoItemssTable,
		listsItemsTable,
		usersListsTable,
		setQuery,
	)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
