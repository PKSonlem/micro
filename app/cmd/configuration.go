package cmd

import (
	"fmt"
)

const (
	envPostgresDB                 = "POSTGRES_DB"
	envPostgresHost               = "POSTGRES_HOST"
	envPostgresPort               = "POSTGRES_PORT"
	envPostgresUser               = "POSTGRES_USER"
	envPostgresPassword           = "POSTGRES_PASSWORD"
	envPostgresSslMode            = "POSTGRES_SSL_MODE"
	envPostgresMaxIdleConnections = "POSTGRES_MAX_IDLE_CONNECTIONS"
	envPostgresMaxOpenConnections = "POSTGRES_MAX_OPEN_CONNECTIONS"

	envServerHost = "SERVER_HOST"
	envServerPort = "SERVER_PORT"

	envJWTSecret = "JWT_SECRET"

	envOutboxWorkerInterval = "OUTBOX_WORKER_INTERVAL"
	envOutboxWorkerLimit    = "OUTBOX_WORKER_LIMIT"
)

// newFromEnv создает конфигурацию и загружает все переменные окружения сразу
func newFromEnv() (*configuration, error) {
	c := &configuration{}

	// Инициализируем PostgreSQL конфигурацию
	pc := &postgresConfiguration{}

	var err error
	pc.user, err = getStringFromEnv(envPostgresUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres user: %w", err)
	}

	pc.host, err = getStringFromEnv(envPostgresHost)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres host: %w", err)
	}

	pc.port, err = getIntValueFromEnv(envPostgresPort, 5432)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres port: %w", err)
	}

	pc.password, err = getStringFromEnv(envPostgresPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres password: %w", err)
	}

	pc.db, err = getStringFromEnv(envPostgresDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres db: %w", err)
	}

	pc.sslmode = getStringFromEnvOrDefault(envPostgresSslMode, "disable")

	pc.maxIdleConnections, err = getIntValueFromEnv(envPostgresMaxIdleConnections, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres max idle connections: %w", err)
	}

	pc.maxOpenConnections, err = getIntValueFromEnv(envPostgresMaxOpenConnections, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres max open connections: %w", err)
	}

	c.postgresConfiguration = pc

	// Инициализируем Server конфигурацию
	sc := &serverConfiguration{}

	sc.host = getStringFromEnvOrDefault(envServerHost, "localhost")

	sc.port, err = getIntValueFromEnv(envServerPort, 8080)
	if err != nil {
		return nil, fmt.Errorf("failed to get server port: %w", err)
	}

	c.serverConfiguration = sc

	// Загружаем JWT секрет
	jwtSecret := getStringFromEnvOrDefault(envJWTSecret, "default-secret-key-change-me")
	c.jwtSecret = jwtSecret

	wc := &workerConfiguration{}

	wc.interval, err = getIntValueFromEnv(envOutboxWorkerInterval, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get outbox worker poll interval: %w", err)
	}

	wc.limit, err = getIntValueFromEnv(envOutboxWorkerLimit, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get outbox worker batch size: %w", err)
	}

	c.workerConfiguration = wc

	return c, nil
}

// структура для хранения конфигураций, под каждую новую зависимость переменные окружения парсятся тут
type configuration struct {
	postgresConfiguration *postgresConfiguration
	serverConfiguration   *serverConfiguration
	workerConfiguration   *workerConfiguration
	jwtSecret             string
}

type postgresConfiguration struct {
	db                 string
	host               string
	port               int64
	user               string
	password           string
	sslmode            string
	maxIdleConnections int64
	maxOpenConnections int64
}

type serverConfiguration struct {
	host string
	port int64
}

type workerConfiguration struct {
	interval int64
	limit    int64
}

func (c *configuration) GetPostgresConfiguration() *postgresConfiguration {
	return c.postgresConfiguration
}

func (pc *postgresConfiguration) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		pc.host,
		pc.port,
		pc.user,
		pc.db,
		pc.password,
		pc.sslmode)
}

func (pc *postgresConfiguration) GetMigrateConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", pc.user, pc.password, pc.host, pc.port, pc.db, pc.sslmode)
}

func (pc *postgresConfiguration) GetMaxIdleConns() int {
	return int(pc.maxIdleConnections)
}

func (pc *postgresConfiguration) GetMaxOpenConns() int {
	return int(pc.maxOpenConnections)
}

func (c *configuration) GetServerConfiguration() *serverConfiguration {
	return c.serverConfiguration
}

func (sc *serverConfiguration) GetAddress() string {
	return fmt.Sprintf("%s:%d", sc.host, sc.port)
}

func (c *configuration) GetJWTSecret() string {
	return c.jwtSecret
}

func (c *configuration) GetWorkerConfiguration() *workerConfiguration {
	return c.workerConfiguration
}

func (w *workerConfiguration) GetInterval() int64 {
	return w.interval
}

func (w *workerConfiguration) GetLimit() int {
	return int(w.limit)
}
