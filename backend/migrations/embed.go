// Package migrations embeds the SQL migration files so the migrate
// runner binary works without filesystem access.
package migrations

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS
