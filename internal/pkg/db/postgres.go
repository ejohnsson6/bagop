package db

import (
	"fmt"
)

// DumpPostgresCmd returns the command to dump a PostgreSQL database given it's environment variables
func DumpPostgresCmd(env []string) []string {

	username := findEnvVar(env, "POSTGRES_USER")
	db := findEnvVar(env, "POSTGRES_DB")

	command := []string{"pg_dump", fmt.Sprintf("--username=%s", username), db}

	return command

}
