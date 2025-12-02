package auth

import (
	"context"
	"fmt"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo         repository
	tokenService tokenService
}

func NewAuthUsecase(repository repository, tokenService tokenService) *AuthUsecase {
	return &AuthUsecase{
		repo:         repository,
		tokenService: tokenService,
	}
}

func (a *AuthUsecase) HandleRegister(ctx context.Context, email, password, userType string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password could not be hashed: %w", err)
	}

	user := entity.User{
		Email:        email,
		PasswordHash: string(hashPassword),
		UserType:     userType,
	}

	cUser, err := a.repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error creared user: %w", err)
	}

	return cUser.UserId, nil
}

func (a *AuthUsecase) HandleLogin(ctx context.Context, userId, password string) (string, error) {
	user, err := a.repo.GetUserById(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("receiving error user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("сouldn't compare passwords: %w", err)
	}

	role := entity.Role(user.UserType)
	token, err := a.tokenService.GenerateToken(userId, role)
	if err != nil {
		return "", fmt.Errorf("error generated token: %w", err)
	}

	return token, nil
}
