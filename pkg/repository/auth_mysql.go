package repository

import (
	"fmt"
	"todo-app/types"

	"github.com/jmoiron/sqlx"
)

type AuthMysql struct {
	db *sqlx.DB
}

func NewAuthMySql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{db: db}
}

func (r *AuthMysql) CreateUser(user types.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values (?, ?, ?)", usersTable)

	result, err := r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *AuthMysql) GetUser(username, password string) (types.User, error) {
	var user types.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=? AND password_hash=?", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
