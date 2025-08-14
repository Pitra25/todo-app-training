package models

import "errors"

const ( // table DB
	UsersTable      = "users"
	TodoListsTable  = "todo_lists"
	UsersListsTable = "users_lists"
	TodoItemsTable  = "todo_items"
	ListsItemsTable = "lists_items"
)

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UserList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItems struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListItems struct {
	Id     int
	ListId int
	ItemId int
}

type UpdadeListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdadeListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdadeItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done" db:"done"`
}

func (i UpdadeItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
