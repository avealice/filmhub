package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/avealice/filmhub/internal/model"
	"github.com/avealice/filmhub/internal/repository"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "djfdfnnvnfbnv"
	signingKey = "dkdjmvfdivs"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Role   string `json:"role"`
	UserID int    `json:"user_id"`
}

type AuthService struct {
	r repository.Authorization
}

func NewAuthService(r repository.Authorization) *AuthService {
	return &AuthService{
		r: r,
	}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.r.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.r.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	claims := &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role:   user.Role,
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return -1, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return -1, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, claims.Role, nil
}
