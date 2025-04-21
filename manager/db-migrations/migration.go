//go:build linux
// +build linux

package dbmigration

import "embed"

//go:embed db/migrations/*.sql
var MigrationsFs embed.FS
