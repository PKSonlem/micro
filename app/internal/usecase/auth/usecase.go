package auth

import (
	"context"
	"fmt"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/pkg/events"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo           repository
	tokenService   tokenService
	eventPublisher *events.Publisher
}

func NewAuthUsecase(repository repository, tokenService tokenService, eventPublisher *events.Publisher) *AuthUsecase {
	return &AuthUsecase{
		repo:           repository,
		tokenService:   tokenService,
		eventPublisher: eventPublisher,
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

	a.eventPublisher.PublishWithUser(ctx, events.EventUserRegistered, 0, "user", map[string]interface{}{
		"user_id":   cUser.UserId,
		"email":     email,
		"user_type": userType,
	})

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

	a.eventPublisher.PublishWithUser(ctx, events.EventUserLogin, 0, "user", map[string]interface{}{
		"user_id": userId,
	})

	return token, nil
}
