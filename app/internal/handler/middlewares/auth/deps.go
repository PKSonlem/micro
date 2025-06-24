package auth

import (
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/service/token"
)

//go:generate mockgen -source=deps.go -destination=mocks/mock_token_service.go -package=mocks

// TokenService интерфейс для работы с токенами
type TokenService interface {
	GenerateToken(userID string, role entity.Role) (string, error)
	ValidateToken(tokenString string) (*token.Claims, error)
}
