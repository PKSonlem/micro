package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	authhandler "github.com/timurzdev/mentorship-test-task/internal/handler/auth"
	flathandler "github.com/timurzdev/mentorship-test-task/internal/handler/flat"
	househandler "github.com/timurzdev/mentorship-test-task/internal/handler/house"

	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/auth"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/prometheus"
	"github.com/timurzdev/mentorship-test-task/internal/handler/server"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
	"github.com/timurzdev/mentorship-test-task/internal/service/roles"
	"github.com/timurzdev/mentorship-test-task/internal/service/token"
	authusecase "github.com/timurzdev/mentorship-test-task/internal/usecase/auth"
	flatusecases "github.com/timurzdev/mentorship-test-task/internal/usecase/flat"
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
	ctx                context.Context
	db                 *sqlx.DB
	postgresContainer  *postgres.PostgresContainer
	testContainerClean func()
	migrator           *migrations.Migrator
	logger             *logger.Logger
	metrics            *metrics.PrometheusMetrics

	// Services
	tokenService  *token.TokenService
	rolesProvider *roles.RolesProvider

	// Repositories
	repository *repository.Repository

	// Use cases
	authUsecase  *authusecase.AuthUsecase
	houseUsecase *houseusecases.Usecase
	flatUsecase  *flatusecases.Usecase

	// Handlers
	houseHandler *househandler.Handler
	flatHandler  *flathandler.Handler
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
		if err = c.db.Close(); err != nil {
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

func (c *Container) GetPostgresTestContainer() (*postgres.PostgresContainer, func(), error) {
	if c.postgresContainer == nil {
		ctx := context.Background()
		pgConfig := c.configuration.GetPostgresConfiguration()

		container, err := postgres.Run(ctx, "postgres:16-alpine",
			postgres.WithDatabase(pgConfig.db),
			postgres.WithUsername(pgConfig.user),
			postgres.WithPassword(pgConfig.password),
			testcontainers.WithWaitStrategy(
				// Ждем сообщение о готовности базы данных (появляется дважды)
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second),
			),
			testcontainers.WithWaitStrategy(
				// Ждем, когда порт станет доступен (важно для Mac/Windows)
				wait.ForListeningPort("5432/tcp").
					WithStartupTimeout(30*time.Second),
			),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
		}

		c.postgresContainer = container
		c.testContainerClean = func() {
			if tErr := container.Terminate(ctx); tErr != nil {
				c.logger.Error(ctx, fmt.Errorf("failed to terminate postgres container: %w", tErr))
			}
		}
	}

	return c.postgresContainer, c.testContainerClean, nil
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
	if c.repository == nil {
		c.repository = repository.NewRepository(c.db)
	}
	return c.repository
}

func (c *Container) GetHouseUsecase() *houseusecases.Usecase {
	if c.houseUsecase == nil {
		c.houseUsecase = houseusecases.NewUsecase(c.GetRepository())
	}
	return c.houseUsecase
}

func (c *Container) GetFlatUsecase() *flatusecases.Usecase {
	if c.flatUsecase == nil {
		c.flatUsecase = flatusecases.NewUsecase(c.GetRepository())
	}
	return c.flatUsecase
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
		c.houseHandler = househandler.NewHandler(
			c.GetRolesProvider(),
			c.GetHouseUsecase(),
			c.logger,
		)
	}
	return c.houseHandler
}

func (c *Container) GetFlatHandler() *flathandler.Handler {
	if c.flatHandler == nil {
		c.flatHandler = flathandler.NewHandler(
			c.GetFlatUsecase(),
			c.logger,
		)
	}
	return c.flatHandler
}

func (c *Container) GetAuthUsecase() *authusecase.AuthUsecase {
	if c.authUsecase == nil {
		c.authUsecase = authusecase.NewAuthUsecase(
			c.GetRepository(),
			c.GetTokenService(),
		)
	}

	return c.authUsecase
}

func (c *Container) GetAuthHandler() *authhandler.Handler {
	if c.authHandler == nil {
		c.authHandler = authhandler.NewHandler(
			c.GetAuthUsecase(),
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
		flatHandler := c.GetFlatHandler()
		authHandler := c.GetAuthHandler()
		c.server = server.NewServer(
			c.logger,
			c.configuration.GetServerConfiguration().GetAddress(),
			houseHandler,
			flatHandler,
			authHandler,
			c.GetPrometheusMiddleware(),
			c.GetAuthMiddleware(),
		)
	}
	return c.server
}

// InitForTesting создает контейнер для тестирования с testcontainers
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
	// Вместо этого используем testcontainers через геттер
	// Тесты сами управляют жизненным циклом контейнера

	closer := func() {
		// В тестах ничего не закрываем, так как тесты сами управляют ресурсами
	}

	return c, closer, nil
}

// GetTestContainerConnectionString возвращает строку подключения к тестовому PostgreSQL контейнеру
func (c *Container) GetTestContainerConnectionString() (string, error) {
	if c.postgresContainer == nil {
		return "", fmt.Errorf("postgres test container is not initialized")
	}
	return c.postgresContainer.ConnectionString(c.ctx, "sslmode=disable")
}

// GetTestContainerMigrateConnectionString возвращает строку подключения для миграций к тестовому контейнеру
func (c *Container) GetTestContainerMigrateConnectionString() (string, error) {
	if c.postgresContainer == nil {
		return "", fmt.Errorf("postgres test container is not initialized")
	}

	// Получаем обычную строку подключения
	connStr, err := c.postgresContainer.ConnectionString(c.ctx, "sslmode=disable")
	if err != nil {
		return "", err
	}

	// Преобразуем её для использования с migrate (требуется URL формат)
	// ConnectionString возвращает формат: postgres://user:password@host:port/database?sslmode=disable
	return connStr, nil
}
