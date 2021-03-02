package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	l "github.com/swexbe/bagop/internal/pkg/logging"
)

const (
	backupLocation    = "/tmp/bagop"
	backupDBLocation  = backupLocation + "/db"
	extraLocation     = "/extra"
	archiveIDLocation = "/var/bagop/ids.log"
	version           = "1.1.0"
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
	version := flag.Bool("v", false, "Version: Display version")

	flag.Parse()
	if *clean {
		cleanBackups()
	} else if *backup {
		makeBackup()
	} else if *version {
		fmt.Printf("bagop v%s", version)
	} else {
		flag.PrintDefaults()
	}

}
