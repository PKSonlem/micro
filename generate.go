//go:build tools
// +build tools

package generate

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package=generated --config=./config.yaml ./api/api.yaml
