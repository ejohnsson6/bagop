package main

import (
	"log"
	"time"

	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/db"
	"github.com/swexbe/bagop/internal/pkg/docker"
	"github.com/swexbe/bagop/internal/pkg/file"
)

const (
	backupLocation = "/backups/"
)

func main() {
	log.Println("Looking for labled containers")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := docker.GetEnabledContainers(cli)
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		vendor, err := docker.FindVendor(container)
		if err != nil {
			panic(err)
		}
		env := docker.GetEnv(cli, container)
		var cmd []string
		containerName := docker.FindName(container)
		switch vendor {
		case "mysql", "mariadb":
			log.Printf("Dumping MYSQL/MariaDB container %s with name %s", container.ID, containerName)
			cmd = db.DumpMysqlCmd(env)
		case "postgres":
			log.Printf("Dumping PostgreSQL container %s with name %s", container.ID, containerName)
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
		timestamp := time.Now().Format(time.RFC3339)
		file.ReaderToFile(reader, dir, timestamp)
	}
}
