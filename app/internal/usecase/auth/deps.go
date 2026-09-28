package auth

import (
	"context"

	"github.com/PKSonlem/micro/internal/entity"
)

//go:generate mockgen -source=deps.go -destination=mock/deps.go -package=mock
type repository interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserById(ctx context.Context, id string) (*entity.User, error)
}

type tokenService interface {
	GenerateToken(userID string, userType entity.Role) (string, error)
}
