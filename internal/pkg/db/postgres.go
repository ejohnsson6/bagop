package db

import (
	"fmt"
)

// DumpPostgresCmd returns the command to dump a PostgreSQL database given it's environment variables
func DumpPostgresCmd(env []string) []string {

	username := findEnvVar(env, "POSTGRES_USER")
	if username == "" {
		// Default username (https://hub.docker.com/_/postgres)
		username = "postgres"
	}
	db := findEnvVar(env, "POSTGRES_DB")
	if db == "" {
		// Default DB name will be set to the value of $POSTGRES_USER (https://hub.docker.com/_/postgres)
		db = username
	}

	command := []string{"pg_dump", fmt.Sprintf("--username=%s", username), db}

	return command

}
