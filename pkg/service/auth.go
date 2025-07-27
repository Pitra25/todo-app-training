package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	"todo-app/pkg/repository"
	"todo-app/types"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt        = "asdjklu48u9r8qwe7244213fw"
	singningKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL    = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenChaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user types.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenChaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(singningKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenChaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid singning method")
		}
		return []byte(singningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenChaims)
	if !ok {
		return 0, errors.New("token chaims are not of type *tokenChaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%s", hash.Sum([]byte(salt)))
}
