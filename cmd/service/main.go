package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/timurzdev/mentorship-test-task/cmd"
)

const (
	codeError = 1
)

// точка входа в нашу программу
func main() {
	//загружаем в окружение переменные из .env файла
	godotenv.Load()

	container := cmd.NewInternal(cmd.NewContainer())

	globalCtx := container.GetGlobalContext()
	log := container.GetLogger()

	ctxFields := map[string]string{
		"path": "cmd/service/main.go",
		"name": "main",
	}

	ctx := log.WithFields(globalCtx, ctxFields)
	log.Info(ctx, "logger initialized")

	migrator := container.GetMigrator()
	err := migrator.MigrateUp()
	if err != nil {
		log.Error(ctx, errors.Wrap(err, "error during migration"))
		os.Exit(codeError)
	}

	log.Info(ctx, "successfull migration")

	server := container.GetServer()
	// TODO: реализовать graceful shutdown
	server.Run(ctx)
}
