package auth

import (
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/service/token"
)

//go:generate mockgen -source=deps.go -destination=mocks/mock_token_service.go -package=mocks

// TokenService интерфейс для работы с токенами
type TokenService interface {
	GenerateToken(userID string, role entity.Role) (string, error)
	ValidateToken(tokenString string) (*token.Claims, error)
}
