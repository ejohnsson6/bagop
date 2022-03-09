package main

import (
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"os"
	"strings"
	"sync"
	"time"

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

func cronJobBackup(ttl string, vaultName string, schedule string, lock *sync.Mutex) {
	l.Logger.Infof("%s: Waiting for lock", schedule)
	lock.Lock()
	defer lock.Unlock()
	l.Logger.Infof("Running %s scheduled backup", schedule)
	makeBackup(ttl, vaultName)
	cleanBackups(vaultName)
	l.Logger.Infof("%s: Letting go of lock", schedule)
}

func main() {
	_ = godotenv.Load()
	cron := flag.Bool("s", false, "Scheduled: run various backups and clean jobs on a schedule based en environment variables")
	clean := flag.Bool("c", false, "Clean: Remove archives which have expired")
	backup := flag.Bool("b", false, "Backup: Make a backup and push it to Glacier")
	list := flag.Bool("l", false, "List: Pretty print all non-expired archives")
	version := flag.Bool("version", false, "Version: Display version")
	verbose := flag.Bool("v", false, "Verbose: Display debug information")
	forceColor := flag.Bool("force-color", false, "Force Color: Force output to be color")
	ttl := flag.String("ttl", "", "Time to Live: Number of days until archives will be deleted")
	vaultName := os.Getenv(utility.ENVVault)
	verboseEnv := os.Getenv(utility.ENVVerbose)
	colorEnv := os.Getenv(utility.ENVColor)

	flag.Parse()

	if *forceColor || (colorEnv != "" && strings.ToLower(colorEnv) != "false") {
		l.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	if *verbose || (verboseEnv != "" && strings.ToLower(verboseEnv) != "false") {
		l.Logger.SetLevel(logrus.DebugLevel)
		l.Logger.Infof("Running in verbose mode")
	}
	if *version {
		fmt.Printf("bagop v%s\n", utility.Version)
		os.Exit(0)
	}
	if *list {
		printArchives()
		os.Exit(0)
	}
	if *clean {
		cleanBackups(vaultName)
	}
	if *backup {
		makeBackup(*ttl, vaultName)
	}
	if *cron {
		cronEnv := os.Getenv(utility.ENVCron)
		ltCronEnv := os.Getenv(utility.ENVLTCron)
		ttlEnv := os.Getenv(utility.ENVTTL)
		ltTtlEnv := os.Getenv(utility.ENVLTTTL)
		s := gocron.NewScheduler(time.Local)
		var m sync.Mutex
		_, err := s.Cron(cronEnv).Do(func() { cronJobBackup(ttlEnv, vaultName, "Regular", &m) })
		panicIfErr(err)
		_, err = s.Cron(ltCronEnv).Do(func() { cronJobBackup(ltTtlEnv, vaultName, "Long-term", &m) })
		panicIfErr(err)
		s.SingletonMode().StartBlocking()
	}
	if !(*version || *clean || *backup || *cron) {
		flag.PrintDefaults()
	}
}
