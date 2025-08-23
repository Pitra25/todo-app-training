package methods

import (
	"fmt"
	"todo-app/internal/repository/mysql/models"

	"github.com/jmoiron/sqlx"
)

type UsersMySQL struct {
	db *sqlx.DB
}

func NewUsersMySQL(db *sqlx.DB) *UsersMySQL {
	return &UsersMySQL{db: db}
}

func (r *UsersMySQL) GetUserById(userId int) (*models.UserResponse, error) {
	var user models.UserResponse
	query := fmt.Sprintf(
		"SELECT id, name, username FROM %s WHERE id = ?",
		models.UsersTable,
	)
	err := r.db.Get(&user, query, userId)
	if err != nil {
		return &models.UserResponse{}, err
	}

	return &user, nil
}


func (r *UsersMySQL) GetUserAll() (*[]models.UserResponse, error) {
	var user []models.UserResponse

	query := fmt.Sprintf(
		"SELECT id, name, username, email FROM %s",
		models.UsersTable,
	)

	err := r.db.Get(&user, query)
	if err != nil {
		return &[]models.UserResponse{}, err
	}

	return &user, nil
}

func (r *UsersMySQL) UpdateUser(userId int, input *models.UpdateUserInput) error {
	query := fmt.Sprintf(
		"UPDATE %s SET name = ?, username = ? WHERE id = ?",
		models.UsersTable,
	)
	_, err := r.db.Exec(query, input.Name, input.Username, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersMySQL) DeleteUser(userId int) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id = ?",
		models.UsersTable,
	)
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}
