package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims структура для JWT токена
type Claims struct {
	UserID   string      `json:"user_id"`
	UserType entity.Role `json:"user_type"`
	jwt.RegisteredClaims
}

type TokenService struct {
	key           []byte
	tokenDuration time.Duration
}

func NewTokenService(key []byte) *TokenService {
	return &TokenService{
		key:           key,
		tokenDuration: 24 * time.Hour, // По умолчанию токен живет 24 часа
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (t *TokenService) GenerateToken(userID string, userType entity.Role) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(t.key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет и парсит JWT токен
func (t *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return t.key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Проверяем срок действия токена
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
