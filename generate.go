//go:build tools
// +build tools

// при запуске go generate ./generate.go мы генерируем модели и хттп сервер на основе openapi схемы, которую мы получили в задании
package generate

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package=generated --config=./config.yaml ./api/api.yaml
