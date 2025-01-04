package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/timurzdev/mentorship-test-task/cmd"
)

const (
	codeError = 1
)

func main() {
	external := cmd.NewContainer()
	internal := cmd.NewInternal(external)

	globalCtx := external.GetGlobalContext()
	log := external.GetLogger()

	ctxFields := map[string]string{
		"path": "cmd/service/main.go",
		"name": "main",
	}

	ctx := log.WithFields(globalCtx, ctxFields)
	log.Info(ctx, "logger initialized")

	migrator := external.GetMigrator()
	err := migrator.MigrateUp()
	if err != nil {
		log.Error(ctx, errors.Wrap(err, "error during migration"))
		os.Exit(codeError)
	}

	log.Info(ctx, "successfull migration")

	server := internal.GetServer()
	// TODO: graceful shutdown
	server.Run(ctx)
}
