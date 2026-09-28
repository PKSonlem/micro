package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	depsmocks "github.com/PKSonlem/micro/internal/deps/mocks"
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/handler/middlewares/auth/mocks"
	"github.com/PKSonlem/micro/internal/service/roles"
	"github.com/PKSonlem/micro/internal/service/token"
	"go.uber.org/mock/gomock"
)

func TestMiddleware_Apply_ExcludedPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	excludedPaths := []string{
		"/dummyLogin",
		"/login",
		"/register",
		"/health",
		"/metrics",
	}

	for _, path := range excludedPaths {
		t.Run("excludes_"+path, func(t *testing.T) {
			// Создаем тестовый handler который записывает статус 200
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Применяем middleware
			handler := middleware.Apply(nextHandler)

			// Создаем запрос без заголовка авторизации
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			// Выполняем запрос
			handler.ServeHTTP(rec, req)

			// Проверяем что запрос прошел успешно без проверки токена
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestMiddleware_Apply_MissingAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.Apply(nextHandler)

	// Создаем запрос на защищенный путь без заголовка
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing authorization header")
}

func TestMiddleware_Apply_InvalidAuthHeaderFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	testCases := []struct {
		name       string
		authHeader string
	}{
		{"no_bearer_prefix", "token123"},
		{"wrong_prefix", "Basic token123"},
		{"only_bearer", "Bearer"},
		{"too_many_parts", "Bearer token part3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.Apply(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.authHeader)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Contains(t, rec.Body.String(), "invalid authorization header format")
		})
	}
}

func TestMiddleware_Apply_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	// Мокаем валидацию токена с ошибкой
	mockTokenService.EXPECT().
		ValidateToken("invalid-token").
		Return(nil, errors.New("invalid token"))

	mockLogger.EXPECT().
		Error(gomock.Any(), gomock.Any())

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.Apply(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid token")
}

func TestMiddleware_Apply_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	// Мокаем валидацию с ошибкой истекшего токена
	mockTokenService.EXPECT().
		ValidateToken("expired-token").
		Return(nil, token.ErrExpiredToken)

	mockLogger.EXPECT().
		Error(gomock.Any(), token.ErrExpiredToken)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.Apply(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "token has expired")
}

func TestMiddleware_Apply_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	claims := &token.Claims{
		UserID:   "user123",
		UserType: entity.RoleClient,
	}

	// Мокаем успешную валидацию токена
	mockTokenService.EXPECT().
		ValidateToken("valid-token").
		Return(claims, nil)

	var capturedCtx context.Context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.Apply(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем что данные добавлены в контекст
	// Используем RolesProvider для проверки
	provider := roles.NewProvider()
	role, err := provider.GetRole(capturedCtx)
	require.NoError(t, err)
	assert.Equal(t, entity.RoleClient, role)

	userID, err := provider.GetUserID(capturedCtx)
	require.NoError(t, err)
	assert.Equal(t, "user123", userID)
}

func TestMiddleware_Apply_ModeratorPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name            string
		path            string
		userRole        entity.Role
		expectedStatus  int
		expectRoleCheck bool
	}{
		{
			name:            "client_forbidden_on_moderator_path",
			path:            "/house/create",
			userRole:        entity.RoleClient,
			expectedStatus:  http.StatusForbidden,
			expectRoleCheck: true,
		},
		{
			name:            "moderator_allowed_on_moderator_path",
			path:            "/house/create",
			userRole:        entity.RoleModerator,
			expectedStatus:  http.StatusOK,
			expectRoleCheck: true,
		},
		{
			name:            "client_allowed_on_regular_path",
			path:            "/some/other/path",
			userRole:        entity.RoleClient,
			expectedStatus:  http.StatusOK,
			expectRoleCheck: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTokenService := mocks.NewMockTokenService(ctrl)
			mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
			mockLogger := depsmocks.NewMockLogger(ctrl)

			middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

			claims := &token.Claims{
				UserID:   "user123",
				UserType: tc.userRole,
			}

			// Мокаем валидацию токена
			mockTokenService.EXPECT().
				ValidateToken("token").
				Return(claims, nil)

			// Мокаем проверку роли если требуется
			if tc.expectRoleCheck {
				mockRolesProvider.EXPECT().
					GetRole(gomock.Any()).
					Return(tc.userRole, nil)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.Apply(nextHandler)

			req := httptest.NewRequest(http.MethodPost, tc.path, nil)
			req.Header.Set("Authorization", "Bearer token")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusForbidden {
				assert.Contains(t, rec.Body.String(), "forbidden: moderator role required")
			}
		})
	}
}

func TestMiddleware_Apply_RoleProviderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mocks.NewMockTokenService(ctrl)
	mockRolesProvider := depsmocks.NewMockRolesProvider(ctrl)
	mockLogger := depsmocks.NewMockLogger(ctrl)

	middleware := NewMiddleware(mockTokenService, mockRolesProvider, mockLogger)

	claims := &token.Claims{
		UserID:   "user123",
		UserType: entity.RoleClient,
	}

	// Мокаем валидацию токена
	mockTokenService.EXPECT().
		ValidateToken("token").
		Return(claims, nil)

	// Мокаем ошибку при получении роли
	roleErr := errors.New("role provider error")
	mockRolesProvider.EXPECT().
		GetRole(gomock.Any()).
		Return(entity.Role(""), roleErr)

	mockLogger.EXPECT().
		Error(gomock.Any(), roleErr)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.Apply(nextHandler)

	req := httptest.NewRequest(http.MethodPost, "/house/create", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}
