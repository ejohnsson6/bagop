package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	version := flag.Bool("version", false, "Version: Display version")
	verbose := flag.Bool("v", false, "Verbose: Display debug information")
	forceColor := flag.Bool("force-color", false, "Force Color: Force output to be color")
	ttl := flag.String("ttl", "", "Time to Live: Number of days until archives will be deleted")
	vaultName := os.Getenv(utility.ENVVault)

	flag.Parse()

	if *forceColor {
		l.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	if *verbose {
		l.Logger.SetLevel(logrus.DebugLevel)
		l.Logger.Infof("Running in verbose mode")
	}
	if *version {
		fmt.Printf("bagop v%s\n", utility.Version)
		os.Exit(0)
	}
	if *clean {
		cleanBackups(vaultName)
	}
	if *backup {
		makeBackup(*ttl, vaultName)
	}
	if !(*version || *clean || *backup) {
		flag.PrintDefaults()
	}
}
