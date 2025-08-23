package methods

import (
	"fmt"
	"todo-app/internal/repository/mysql/models"

	"github.com/jmoiron/sqlx"
)

type AuthMysql struct {
	db *sqlx.DB
}

func NewAuthMySql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{db: db}
}

func (r *AuthMysql) CreateUser(user models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash, email) values (?, ?, ?, ?)", models.UsersTable)

	result, err := r.db.Exec(query, user.Name, user.Username, user.Password, user.Email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *AuthMysql) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=? AND password_hash=?", models.UsersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
