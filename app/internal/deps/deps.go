package deps

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

//go:generate mockgen -source deps.go -destination=mocks/deps.go -package mocks

type Logger interface {
	Info(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, err error, args ...any)
}

type Metrics interface {
	RegisterCounter(key string, c *prometheus.CounterVec) error
	RegisterHistogram(key string, c *prometheus.HistogramVec) error
	GetCounter(key string) (*prometheus.CounterVec, error)
	GetHistogram(key string) (*prometheus.HistogramVec, error)
}

// RolesProvider интерфейс для работы с ролями пользователей
type RolesProvider interface {
	GetRole(ctx context.Context) (entity.Role, error)
	GetUserID(ctx context.Context) (string, error)
}
