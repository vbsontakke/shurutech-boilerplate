// Package migrations embeds the SQL migration files into the binary so they
// run automatically on startup and in tests.
//
// Naming: 0001_init.up.sql, 0002_create_xxx.up.sql, ... (applied in order).
// Matching *.down.sql files are kept for documentation / manual rollback.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
