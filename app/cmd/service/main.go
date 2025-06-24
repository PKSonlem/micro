package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/timurzdev/mentorship-test-task/cmd"
)

const (
	codeError = 1
)

// точка входа в нашу программу
func main() {
	// загружаем в окружение переменные из .env файла
	_ = godotenv.Load("./deploy/local/.env")

	// Инициализируем контейнер
	container, closer, err := cmd.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize container: %v\n", err)
		os.Exit(codeError)
	}
	defer closer()

	log := container.GetLogger()
	ctx := container.GetContext()

	ctxFields := map[string]string{
		"path": "cmd/service/main.go",
		"name": "main",
	}

	ctx = log.WithFields(ctx, ctxFields)
	log.Info(ctx, "application starting")

	// Выполняем миграции
	migrator := container.GetMigrator()
	err = migrator.MigrateUp()
	if err != nil {
		log.Error(ctx, errors.Wrap(err, "error during migration"))
		os.Exit(codeError)
	}

	log.Info(ctx, "successful migration")

	// Запускаем сервер в отдельной горутине
	server := container.GetServer()
	go func() {
		log.Info(ctx, fmt.Sprintf("starting server on %s", container.GetServerAddress()))
		server.Run(ctx)
	}()

	// Настраиваем graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info(ctx, "shutting down server...")

	// Даем серверу время на завершение активных соединений
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: добавить graceful shutdown для HTTP сервера
	<-shutdownCtx.Done()

	log.Info(ctx, "server stopped")
}
