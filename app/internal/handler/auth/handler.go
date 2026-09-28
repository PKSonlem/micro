package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/PKSonlem/micro/internal/deps"
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/generated"
	"github.com/PKSonlem/micro/internal/service/converters"
	"github.com/PKSonlem/micro/internal/service/token"
	"github.com/PKSonlem/micro/internal/usecase/auth"
)

type Handler struct {
	authUsecase  *auth.AuthUsecase
	tokenService *token.TokenService
	logger       deps.Logger
}

func NewHandler(authUsecase *auth.AuthUsecase, tokenService *token.TokenService, logger deps.Logger) *Handler {
	return &Handler{
		authUsecase:  authUsecase,
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

func (h *Handler) PostRegister(ctx context.Context, request generated.PostRegisterRequestObject) (generated.PostRegisterResponseObject, error) {
	if request.Body == nil {
		return generated.PostRegister400Response{}, nil
	}

	user := converters.UserFromGenRegister(generated.PostRegisterJSONBody(*request.Body))

	email := user.Email
	password := user.PasswordHash
	userType := user.UserType

	userId, err := h.authUsecase.HandleRegister(ctx, email, password, userType)
	if err != nil {
		h.logger.Error(ctx, err)
		if errors.Is(err, entity.ErrorCreatingUser) {
			return generated.PostRegister400Response{}, nil
		}

		return generated.PostRegister500JSONResponse{
			N5xxJSONResponse: generated.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Internal server error",
				},
			},
		}, nil
	}

	uuidParsed, err := uuid.Parse(userId)
	if err != nil {
		return generated.PostRegister400Response{}, nil
	}

	userIdGenerated := generated.UserId(uuidParsed)

	return generated.PostRegister200JSONResponse{
		UserId: &userIdGenerated,
	}, nil
}

func (h *Handler) PostLogin(ctx context.Context, request generated.PostLoginRequestObject) (generated.PostLoginResponseObject, error) {
	if request.Body == nil {
		return generated.PostLogin400Response{}, nil
	}

	userId, password := converters.UserLoginFromGen(generated.PostLoginJSONBody(*request.Body))

	token, err := h.authUsecase.HandleLogin(ctx, userId, password)
	if err != nil {
		h.logger.Error(ctx, err)
		if errors.Is(err, entity.ErrorLoginUser) {
			return generated.PostLogin404Response{}, nil
		}
		return generated.PostLogin500JSONResponse{
			N5xxJSONResponse: generated.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Internal server error",
				},
			},
		}, nil
	}

	return generated.PostLogin200JSONResponse{
		Token: &token,
	}, nil
}
