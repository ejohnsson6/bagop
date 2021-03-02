package main

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	l "github.com/swexbe/bagop/internal/pkg/logging"
	"github.com/swexbe/bagop/internal/pkg/utility"
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
	ttl := flag.String("ttl", "", "Time to Live: Number of days until archives will be deleted")

	flag.Parse()
	if *clean {
		cleanBackups()
	} else if *backup {
		makeBackup(*ttl)
	} else if *version {
		fmt.Printf("bagop v%s", utility.Version)
	} else {
		flag.PrintDefaults()
	}
}
