package deps

import (
	"context"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
)

type Logger interface {
	Info(ctx context.Context, message string, args ...any)
	Error(ctx context.Context, err error, args ...any)
}

// не имплементировано, нужно реализовать позже
type RolesProvider interface {
	RolesReader
	RolesWriter
}

type RolesReader interface {
	GetRole(ctx context.Context) (entity.Role, error)
}

type RolesWriter interface {
	SetRole(ctx context.Context, role entity.Role) error
}
