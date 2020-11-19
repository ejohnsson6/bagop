package db

import (
	"fmt"
)

// DumpMysqlCmd returns the command to dump a MYSQL/MariaDB database given it's environment variables
func DumpMysqlCmd(env []string) []string {

	username := findEnvVar(env, "MYSQL_USER")
	password := findEnvVar(env, "MYSQL_PASSWORD")
	db := findEnvVar(env, "MYSQL_DATABASE")

	command := []string{"mysqldump", fmt.Sprintf("-u%s", username), fmt.Sprintf("-p%s", password), "--skip-comments", "--databases", db}

	return command
}
