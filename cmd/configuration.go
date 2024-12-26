package cmd

import (
	"fmt"
	"log"
)

const (
	envPostgresHost               = "POSTGRES_HOST"
	envPosrgresPort               = "POSTGRES_PORT"
	envPostgresUser               = "POSTGRES_USER"
	envPostgresPassword           = "POSTGRES_PASSWORD"
	envPostgresSslMode            = "POSTGRES_SSL_MODE"
	envPostgresMaxIdleConnections = "POSTGRES_MAX_IDLE_CONNECTIONS"
	envPostgresMaxOpenConnections = "POSTGRES_MAX_OPEN_CONNECTIONS"
)

func newFromEnv() *configuration {
	c := &configuration{}

	return c
}

// структура для хранения конфигураций, под каждую новую зависимость переменные окружения парсятся тут
type configuration struct {
	postgresConfgiration *postgresConfiguration
}

type postgresConfiguration struct {
	host               string
	port               int64
	user               string
	password           string
	sslmode            string
	maxIdleConnections int64
	maxOpenConnections int64
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

		pc.host, err = getStringFromEnv(envPostgresUser)
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

func (pc *postgresConfiguration) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s", pc.host, pc.port, pc.user, pc.password, pc.sslmode)
}

func (pc *postgresConfiguration) GetMaxIdleConns() int {
	return int(pc.maxIdleConnections)
}

func (pc *postgresConfiguration) GetMaxOpenConns() int {
	return int(pc.maxOpenConnections)
}
