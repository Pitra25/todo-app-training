package methods

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
	salt           = "asdjklu48u9r8qwe7244213fw"
	signingKey_12 = os.Getenv("SIGNINGKEY_12")
	// signingKey_12 = "qrkjk#4#35FSFJlja#4353KSFjH"
	signingKey_35 = os.Getenv("SIGNINGKEY_35")
	// tokenTTL_A_05  = time.Duration(viper.GetInt64("jwt.timeOfLife_1")) * time.Hour // 12 hours, ccess token
	// tokenTTL_R_35  = time.Duration(viper.GetInt64("jwt.timeOfLife_2")) * time.Hour // 35 days, refresh token
)

type AuthService struct {
	repo repository.Authorization
}

type tokenChains struct {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenChains{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL_A_05).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey_12))
}

func (s *AuthService) GenerateRefreshToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenChains{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL_R_35).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey_35))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenChains{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey_12), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenChains)
	if !ok {
		return 0, errors.New("token chains are not of type *tokenChains")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return string(hash.Sum([]byte(salt)))
}
