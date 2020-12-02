package main

import (
	"time"

	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/aws"
	"github.com/swexbe/bagop/internal/pkg/db"
	"github.com/swexbe/bagop/internal/pkg/docker"
	"github.com/swexbe/bagop/internal/pkg/file"
	l "github.com/swexbe/bagop/internal/pkg/logging"
)

const (
	backupLocation = "/backups/"
)

func main() {
	l.Logger.Infof("Looking for labled containers")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := docker.GetEnabledContainers(cli)
	if err != nil {
		panic(err)
	}
	l.Logger.Infof("Found %d", len(containers))
	timestamp := time.Now().Format(time.RFC3339)

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
		if err != nil {
			panic(err)
		}
		reader, err := docker.RunCommand(cli, container, cmd)
		if err != nil {
			panic(err)
		}

		dir := backupLocation + containerName + "/"
		file.ReaderToFile(reader, dir, timestamp)
	}
	aws.Test()
}
