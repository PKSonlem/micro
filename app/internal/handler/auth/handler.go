package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/service/token"
)

type Handler struct {
	tokenService *token.TokenService
	logger       deps.Logger
}

func NewHandler(tokenService *token.TokenService, logger deps.Logger) *Handler {
	return &Handler{
		tokenService: tokenService,
		logger:       logger,
	}
}

// GetDummyLogin генерирует токен для тестовой авторизации
func (h *Handler) GetDummyLogin(ctx context.Context, request generated.GetDummyLoginRequestObject) (generated.GetDummyLoginResponseObject, error) {
	// Преобразуем UserType из generated в entity.Role
	var role entity.Role
	switch request.Params.UserType {
	case generated.Client:
		role = entity.RoleClient
	case generated.Moderator:
		role = entity.RoleModerator
	default:
		return generated.GetDummyLogin500JSONResponse{
			N5xxJSONResponse: generated.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Invalid user type",
				},
			},
		}, nil
	}

	// Генерируем случайный user ID для dummy login
	userID := uuid.New().String()

	// Генерируем токен
	tokenString, err := h.tokenService.GenerateToken(userID, role)
	if err != nil {
		h.logger.Error(ctx, err)
		return generated.GetDummyLogin500JSONResponse{
			N5xxJSONResponse: generated.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Failed to generate token",
				},
			},
		}, nil
	}

	h.logger.Info(ctx, "dummy login successful", "user_id", userID, "user_type", role)

	return generated.GetDummyLogin200JSONResponse{
		Token: &tokenString,
	}, nil
}
