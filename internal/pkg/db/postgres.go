package db

import (
	"fmt"

	l "github.com/swexbe/bagop/internal/pkg/logging"
)

// DumpPostgresCmd returns the command to dump a PostgreSQL database given it's environment variables
func DumpPostgresCmd(env []string) []string {

	username := findEnvVar(env, "POSTGRES_USER")
	if username == "" {
		// Default username (https://hub.docker.com/_/postgres)
		username = "postgres"
		l.Logger.Warnf("POSTGRES_USER label not found, using postgres default")
	}
	db := findEnvVar(env, "POSTGRES_DB")
	if db == "" {
		// Default DB name will be set to the value of $POSTGRES_USER (https://hub.docker.com/_/postgres)
		db = username
		l.Logger.Warnf("POSTGRES_DB label not found, using postgres default")
	}

	command := []string{"pg_dump", fmt.Sprintf("--username=%s", username), db}

	return command

}
