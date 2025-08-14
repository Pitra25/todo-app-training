package models

import "errors"

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (u User) Validate() error {
	if u.Name == "" && u.Username == "" && u.Password == "" && u.Email == "" {
		return errors.New("update structure has no values")
	}
	return nil
}

type UserResponse struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateUserInpur struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (u UpdateUserInpur) Validate() error {
	if u.Name == nil && u.Username == nil && u.Password == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
