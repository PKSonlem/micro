package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	embPg "github.com/fergusstrange/embedded-postgres"
)

const (
	envPostgresDB                 = "POSTGRES_DB"
	envPostgresHost               = "POSTGRES_HOST"
	envPosrgresPort               = "POSTGRES_PORT"
	envPostgresUser               = "POSTGRES_USER"
	envPostgresPassword           = "POSTGRES_PASSWORD"
	envPostgresSslMode            = "POSTGRES_SSL_MODE"
	envPostgresMaxIdleConnections = "POSTGRES_MAX_IDLE_CONNECTIONS"
	envPostgresMaxOpenConnections = "POSTGRES_MAX_OPEN_CONNECTIONS"

	envServerHost = "SERVER_HOST"
	envServerPort = "SERVER_PORT"
)

func newFromEnv() *configuration {
	c := &configuration{}

	return c
}

// структура для хранения конфигураций, под каждую новую зависимость переменные окружения парсятся тут
type configuration struct {
	postgresConfgiration *postgresConfiguration
	serverConfiguration  *serverConfiguration
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

func (c *configuration) GetPostgresConfiguration() *postgresConfiguration {
	if c.postgresConfgiration == nil {
		var err error
		pc := &postgresConfiguration{}
		c.postgresConfgiration = pc

		pc.user, err = getStringFromEnv(envPostgresUser)
		if err != nil {
			log.Fatal(err)
		}

		pc.host, err = getStringFromEnv(envPostgresHost)
		if err != nil {
			log.Fatal(err)
		}

		pc.port, err = getIntValueFromEnv(envPosrgresPort, 5432)
		if err != nil {
			log.Fatal(err)
		}

		pc.password, err = getStringFromEnv(envPostgresPassword)
		if err != nil {
			log.Fatal(err)
		}

		pc.sslmode = getStringFromEnvOrDefault(envPostgresSslMode, "disable")

		pc.db, err = getStringFromEnv(envPostgresDB)
		if err != nil {
			log.Fatal(err)
		}

		pc.maxIdleConnections, err = getIntValueFromEnv(envPostgresMaxIdleConnections, 10)
		if err != nil {
			log.Fatal(err)
		}

		pc.maxOpenConnections, err = getIntValueFromEnv(envPostgresMaxOpenConnections, 10)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.postgresConfgiration
}

func (pc *postgresConfiguration) GetEmbeddedPostgresConfig() embPg.Config {
	return embPg.Config{}.
		Database(pc.db).
		Username(pc.user).
		Password(pc.password).
		Port(uint32(pc.port)).
		Version(embPg.V16).
		StartTimeout(time.Second * 15).
		Logger(os.Stdout).
		BinaryRepositoryURL("https://repo1.maven.org/maven2")
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
	if c.serverConfiguration == nil {
		var err error
		sc := &serverConfiguration{}
		c.serverConfiguration = sc

		sc.host = getStringFromEnvOrDefault(envServerHost, "localhost")

		sc.port, err = getIntValueFromEnv(envServerPort, 8080)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.serverConfiguration
}

func (sc *serverConfiguration) GetAddress() string {
	return fmt.Sprintf("%s:%d", sc.host, sc.port)
}
