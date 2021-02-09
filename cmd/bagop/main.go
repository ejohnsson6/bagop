package main

import (
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"github.com/swexbe/bagop/internal/pkg/aws"
	"github.com/swexbe/bagop/internal/pkg/db"
	"github.com/swexbe/bagop/internal/pkg/docker"
	"github.com/swexbe/bagop/internal/pkg/file"
	l "github.com/swexbe/bagop/internal/pkg/logging"
)

const (
	backupLocation = "/tmp/bagop/"
)

func panicIfErr(err error) {
	if err != nil {
		l.Logger.Fatalf(err.Error())
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		l.Logger.Infof(err.Error())
	}
	l.Logger.Infof("Looking for labelled containers")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	panicIfErr(err)
	containers, err := docker.GetEnabledContainers(cli)
	panicIfErr(err)
	l.Logger.Infof("Found %d", len(containers))
	timestamp := time.Now().Format(time.RFC3339)

	os.RemoveAll(backupLocation)

	for _, container := range containers {
		l.Logger.Infof("Trying to dump container %s", container.ID[0:12])
		vendor, err := docker.FindVendor(container)
		if err != nil {
			l.Logger.Errorf(err.Error())
			continue
		}
		env := docker.GetEnv(cli, container)
		var cmd []string
		containerName := docker.FindName(container)
		switch vendor {
		case "mysql", "mariadb":
			l.Logger.Infof("Dumping as MYSQL/MariaDB container with name %s", containerName)
			cmd = db.DumpMysqlCmd(env)
		case "postgres":
			l.Logger.Infof("Dumping as PostgreSQL container with name %s", containerName)
			cmd = db.DumpPostgresCmd(env)
		}
		panicIfErr(err)
		reader, err := docker.RunCommand(cli, container, cmd)
		panicIfErr(err)

		file.ReaderToFile(reader, backupLocation, containerName+".sql")

	}
	tarFileLocation := backupLocation + timestamp + ".tar.gz"
	file.FolderToTarGZ(backupLocation, tarFileLocation)

	err = aws.UploadFile(tarFileLocation, timestamp)
	panicIfErr(err)
}
