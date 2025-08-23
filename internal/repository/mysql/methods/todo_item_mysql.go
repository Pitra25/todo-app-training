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

type TodoItemsMySql struct {
	db  *sqlx.DB
	rdb *storage.RedisCLientDB
}

func NewTodoItemsMySql(db *sqlx.DB, rdb *redis.Client) *TodoItemsMySql {
	return &TodoItemsMySql{
		db:  db,
		rdb: &storage.RedisCLientDB{Db: rdb},
	}
}

func (r *TodoItemsMySql) Create(listId int, item models.TodoItems) (int, error) {
	logrus.Debug("list id:", listId, "item:", item)

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES (?, ?)",
		models.TodoItemsTable,
	)
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

	createListItemsQuery := fmt.Sprintf(
		"INSERT INTO %s (list_id, item_id) VALUES (?, ?)",
		models.ListsItemsTable,
	)
	_, err = tx.Exec(createListItemsQuery, listId, itemid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(itemid), tx.Commit()
}

func (r *TodoItemsMySql) GetAllItemsList(userId, listId int) ([]models.TodoItems, error) {
	var items []models.TodoItems

	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE li.list_id = ? AND ul.user_id = ?",
		models.TodoItemsTable, models.ListsItemsTable, models.UsersListsTable,
	)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemsMySql) GetAllItem() ([]models.TodoItems, error) {
	var items []models.TodoItems

	query := fmt.Sprintf(
		"SELECT id, title, description, done FROM %s",
		models.TodoItemsTable,
	)
	if err := r.db.Select(&items, query); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemsMySql) GetById(userId, itemId int) (models.TodoItems, error) {
	var item models.TodoItems

	// Check Redis cache first
	storageRecords, err := r.rdb.Get(itemId, storage.Item)
	if err != nil {
		logrus.Error("method: GetById item.", err.Error())
	} else if storageRecords.Items.Title != "" && storageRecords != nil {
		return storageRecords.Items, nil
	}

	query := fmt.Sprintf(
		"SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti "+
			"INNER JOIN %s li ON li.item_id = ti.id "+
			"INNER JOIN %s ul ON ul.list_id = li.list_id "+
			"WHERE ti.id = ? AND ul.user_id = ?",
		models.TodoItemsTable, models.ListsItemsTable, models.UsersListsTable,
	)
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	// Save to Redis cache
	if err := r.rdb.Create(&storage.Recording{
		ID:    itemId,
		Items: item,
	}); err != nil {
		logrus.Error("err create record redis.", err.Error())
		return item, nil
	}

	return item, nil
}

func (r *TodoItemsMySql) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		"DELETE ti FROM %s ti "+
			"INNER JOIN %s li ON ti.id = li.item_id "+
			"INNER JOIN %s ul ON li.list_id = ul.list_id "+
			"WHERE ul.user_id = ? AND ti.id = ?",
		models.TodoItemsTable, models.ListsItemsTable, models.UsersListsTable,
	)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *TodoItemsMySql) Update(userId, itemId int, input models.UpdateItemInput) error {
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
		models.TodoItemsTable, models.ListsItemsTable, models.UsersListsTable,
		setQuery,
	)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
