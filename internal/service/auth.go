package service

import (
	"crypto/sha1"
	"errors"
	"os"
	"time"
	"todo-app/internal/repository"
	"todo-app/internal/repository/mysql/models"

	"github.com/dgrijalva/jwt-go"
)

const (
	tokenTTL_A_05 = 12 * time.Hour  // 12 hours, ccess token
	tokenTTL_R_35 = 168 * time.Hour // 35 days, refresh token
)

var (
	// salt = os.Getenv("SALT")
	salt = "asdjklu48u9r8qwe7244213fw"
	singningKey_12 = os.Getenv("SINGNINGKEY_12")
	// singningKey_12 = "qrkjk#4#35FSFJlja#4353KSFjH"
	singningKey_35 = os.Getenv("SINGNINGKEY_35")
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

func (s *AuthService) CreateUser(user models.User) (int, error) {
	if err := user.Validate(); err != nil {
		return 0, err
	}

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
			ExpiresAt: time.Now().Add(tokenTTL_A_05).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(singningKey_12))
}

func (s *AuthService) GenerateRefrachToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenChaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL_R_35).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(singningKey_35))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenChaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid singning method")
		}
		return []byte(singningKey_12), nil
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

	return string(hash.Sum([]byte(salt)))
}
