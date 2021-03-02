package main

import (
	"flag"

	"github.com/joho/godotenv"
	l "github.com/swexbe/bagop/internal/pkg/logging"
)

const (
	backupLocation    = "/tmp/bagop"
	backupDBLocation  = backupLocation + "/db"
	extraLocation     = "/extra"
	archiveIDLocation = "/var/bagop/ids.log"
)

func panicIfErr(err error) {
	if err != nil {
		l.Logger.Fatalf(err.Error())
	}
}

func main() {
	godotenv.Load()
	clean := flag.Bool("c", false, "Clean: Remove archives which have expired")
	backup := flag.Bool("b", false, "Backup: Make a backup and push it to Glacier")

	flag.Parse()
	if *clean {
		cleanBackups()
	} else if *backup {
		makeBackup()
	} else {
		flag.PrintDefaults()
	}

}
