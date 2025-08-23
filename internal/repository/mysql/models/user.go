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
		return errors.New("create structure has no values")
	}
	return nil
}

type UserResponse struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (u UpdateUserInput) Validate() error {
	if u.Name == nil && u.Username == nil && u.Password == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UsersCode struct {
	Id        int    `json:"id" db:"id"`
	UserId    int    `json:"user_id" db:"user_id"`
	Code      string `json:"code" db:"code"`
	ExpiresAt string `json:"expires_at" db:"expires_at"`
}
	
func (uc UsersCode) Validate() error {
	if uc.UserId <= 0 || uc.Code == "" || uc.ExpiresAt == "" {
		return errors.New("invalid user code data")
	}
	return nil
}

type SendMassEmailInput struct {
	To []string `json:"to" binding:"required"`
}

type SendEmailInput struct {
	To string `json:"to" binding:"required"`
}

type SendConfirmationEmailInput struct {
	Code string `json:"code" binding:"required"`
}
