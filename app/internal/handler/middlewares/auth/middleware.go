package auth

import (
	"net/http"
	"strings"

	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/service/roles"
	"github.com/timurzdev/mentorship-test-task/internal/service/token"
)

type Middleware struct {
	tokenService   TokenService
	rolesProvider  deps.RolesProvider
	logger         deps.Logger
	excludedPaths  map[string]bool
	moderatorPaths map[string]bool
}

func NewMiddleware(
	tokenService TokenService,
	rolesProvider deps.RolesProvider,
	logger deps.Logger,
) *Middleware {
	// Пути, которые не требуют авторизации
	excludedPaths := map[string]bool{
		"/dummyLogin": true,
		"/login":      true,
		"/register":   true,
		"/health":     true,
		"/metrics":    true,
	}

	// Пути, требующие роль модератора
	moderatorPaths := map[string]bool{
		"/house/create": true,
		"/flat/update":  true,
	}

	return &Middleware{
		tokenService:   tokenService,
		rolesProvider:  rolesProvider,
		logger:         logger,
		excludedPaths:  excludedPaths,
		moderatorPaths: moderatorPaths,
	}
}

// Apply реализует middleware интерфейс func(http.Handler) http.Handler
func (m *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, исключен ли путь из авторизации
		if m.excludedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Проверяем формат Bearer токена
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Валидируем токен
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Error(r.Context(), err)
			if err == token.ErrExpiredToken {
				http.Error(w, "token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем роль и user_id в контекст
		ctx := roles.SetRoleAndUserID(r.Context(), claims.UserType, claims.UserID)
		r = r.WithContext(ctx)

		// Проверяем, требуется ли роль модератора для этого пути
		if m.moderatorPaths[r.URL.Path] {
			var role entity.Role
			role, err = m.rolesProvider.GetRole(r.Context())
			if err != nil {
				m.logger.Error(r.Context(), err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !role.IsModerator() {
				http.Error(w, "forbidden: moderator role required", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
