package main

import (
	"github.com/timurzdev/mentorship-test-task/cmd"
)

func main() {
	external := cmd.NewContainer()
	internal := cmd.NewInternal()

	globalCtx := external.GetGlobalContext()
	log := external.GetLogger()

	ctx := log.WithTag(globalCtx, "cmd/service/main.go")
	log.Info(ctx, "logger initialized")

	_ = internal.GetRepository(external.GetPostgres())
}
