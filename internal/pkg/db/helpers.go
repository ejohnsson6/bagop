package db

import "strings"

func findEnvVar(env []string, find string) string {
	for _, e := range env {
		split := strings.Split(e, "=")
		if split[0] == find {
			return split[1]
		}
	}
	return ""
}
