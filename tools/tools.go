//go:build tools

package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
)

//go:generate go build -o ../bin/migrate github.com/golang-migrate/migrate/v4
