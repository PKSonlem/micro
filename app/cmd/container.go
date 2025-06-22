package cmd

import (
	"context"
	"fmt"

	embPg "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	authhandler "github.com/timurzdev/mentorship-test-task/internal/handler/auth"
	househandler "github.com/timurzdev/mentorship-test-task/internal/handler/house"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/auth"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/prometheus"
	"github.com/timurzdev/mentorship-test-task/internal/handler/server"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
	"github.com/timurzdev/mentorship-test-task/internal/service/roles"
	"github.com/timurzdev/mentorship-test-task/internal/service/token"
	houseusecases "github.com/timurzdev/mentorship-test-task/internal/usecase/house"
	"github.com/timurzdev/mentorship-test-task/migrations"
	"github.com/timurzdev/mentorship-test-task/pkg/logger"
	"github.com/timurzdev/mentorship-test-task/pkg/metrics"
)

// Container - единый контейнер всех зависимостей приложения
type Container struct {
	// Configuration
	configuration *configuration

	// Infrastructure
	ctx              context.Context
	db               *sqlx.DB
	embeddedPostgres *embPg.EmbeddedPostgres
	migrator         *migrations.Migrator
	logger           *logger.Logger
	metrics          *metrics.PrometheusMetrics

	// Services
	tokenService  *token.TokenService
	rolesProvider *roles.RolesProvider

	// Repositories
	repository *repository.Repository

	// Use cases
	houseUsecase *houseusecases.Usecase

	// Handlers
	houseHandler *househandler.Handler
	authHandler  *authhandler.Handler

	// Server & Middleware
	server               *server.Server
	prometheusMiddleware *prometheus.Middleware
	authMiddleware       *auth.Middleware
}

// Init инициализирует только критические зависимости и возвращает контейнер, closer функцию и ошибку
func Init() (*Container, func(), error) {
	c := &Container{}

	config, err := newFromEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	c.configuration = config

	c.ctx = context.Background()

	c.logger = logger.New()

	c.metrics = metrics.Init()

	c.db, err = NewSqlxConn(c.configuration.GetPostgresConfiguration())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	closer := func() {
		if err := c.db.Close(); err != nil {
			c.logger.Error(c.ctx, fmt.Errorf("failed to close database: %w", err))
		}
	}

	return c, closer, nil
}

// Getters с ленивой инициализацией

// GetServerAddress возвращает адрес сервера
func (c *Container) GetServerAddress() string {
	return c.configuration.GetServerConfiguration().GetAddress()
}

// GetPostgresConnectionString возвращает строку подключения к PostgreSQL
func (c *Container) GetPostgresConnectionString() string {
	return c.configuration.GetPostgresConfiguration().GetConnectionString()
}

func (c *Container) GetContext() context.Context {
	return c.ctx
}

func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Container) GetMetrics() *metrics.PrometheusMetrics {
	return c.metrics
}

func (c *Container) GetDB() *sqlx.DB {
	return c.db
}

func (c *Container) GetEmbeddedPostgres() *embPg.EmbeddedPostgres {
	if c.embeddedPostgres == nil {
		c.embeddedPostgres = embPg.NewDatabase(
			c.configuration.GetPostgresConfiguration().GetEmbeddedPostgresConfig(),
		)
	}
	return c.embeddedPostgres
}

func (c *Container) GetMigrator() *migrations.Migrator {
	if c.migrator == nil {
		c.migrator = migrations.NewMigrator(
			c.configuration.GetPostgresConfiguration().GetMigrateConnectionString(),
		)
	}
	return c.migrator
}

func (c *Container) GetRepository() *repository.Repository {
	if c.repository == nil && c.db != nil {
		c.repository = repository.NewRepository(c.db)
	}
	return c.repository
}

func (c *Container) GetHouseUsecase() *houseusecases.Usecase {
	if c.houseUsecase == nil {
		repo := c.GetRepository()
		if repo != nil {
			c.houseUsecase = houseusecases.NewUsecase(repo)
		}
	}
	return c.houseUsecase
}

func (c *Container) GetTokenService() *token.TokenService {
	if c.tokenService == nil {
		c.tokenService = token.NewTokenService([]byte(c.configuration.GetJWTSecret()))
	}
	return c.tokenService
}

func (c *Container) GetRolesProvider() *roles.RolesProvider {
	if c.rolesProvider == nil {
		c.rolesProvider = roles.NewProvider()
	}
	return c.rolesProvider
}

func (c *Container) GetHouseHandler() *househandler.Handler {
	if c.houseHandler == nil {
		usecase := c.GetHouseUsecase()
		if usecase != nil {
			c.houseHandler = househandler.NewHandler(
				usecase,
				c.logger,
			)
		}
	}
	return c.houseHandler
}

func (c *Container) GetAuthHandler() *authhandler.Handler {
	if c.authHandler == nil {
		c.authHandler = authhandler.NewHandler(
			c.GetTokenService(),
			c.logger,
		)
	}
	return c.authHandler
}

func (c *Container) GetAuthMiddleware() *auth.Middleware {
	if c.authMiddleware == nil {
		c.authMiddleware = auth.NewMiddleware(
			c.GetTokenService(),
			c.GetRolesProvider(),
			c.logger,
		)
	}
	return c.authMiddleware
}

func (c *Container) GetPrometheusMiddleware() *prometheus.Middleware {
	if c.prometheusMiddleware == nil {
		c.prometheusMiddleware = prometheus.New(c.metrics)
	}
	return c.prometheusMiddleware
}

func (c *Container) GetServer() *server.Server {
	if c.server == nil {
		houseHandler := c.GetHouseHandler()
		authHandler := c.GetAuthHandler()
		if houseHandler != nil && authHandler != nil {
			c.server = server.NewServer(
				c.logger,
				c.configuration.GetServerConfiguration().GetAddress(),
				houseHandler,
				authHandler,
				c.GetPrometheusMiddleware(),
				c.GetAuthMiddleware(),
			)
		}
	}
	return c.server
}

// InitForTesting создает контейнер для тестирования с embedded postgres
func InitForTesting() (*Container, func(), error) {
	c := &Container{}

	config, err := newFromEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	c.configuration = config

	c.ctx = context.Background()
	c.logger = logger.New()
	c.metrics = metrics.Init()

	// Для тестов не подключаемся к реальной БД
	// Вместо этого используем embedded postgres через геттер
	// Тесты сами управляют жизненным циклом embedded postgres (Start/Stop)

	closer := func() {
		// В тестах ничего не закрываем, так как тесты сами управляют ресурсами
	}

	return c, closer, nil
}
