package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/aws"
	"github.com/swexbe/bagop/internal/pkg/db"
	"github.com/swexbe/bagop/internal/pkg/docker"
	"github.com/swexbe/bagop/internal/pkg/file"
	l "github.com/swexbe/bagop/internal/pkg/logging"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

func makeBackup() {
	l.Logger.Infof("Looking for labelled containers")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	panicIfErr(err)
	containers, err := docker.GetEnabledContainers(cli)
	panicIfErr(err)
	l.Logger.Infof("Found %d", len(containers))

	if len(containers) == 0 {
		l.Logger.Warnf("No labelled containers found, exiting...")
		return
	}

	timestamp := time.Now()
	timestampStr := timestamp.Format(time.RFC3339)

	os.RemoveAll(utility.BackupLocation)

	for _, container := range containers {
		l.Logger.Infof("Trying to dump container %s", container.ID[0:12])
		vendor, err := docker.FindVendor(cli, container)
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

		file.ReaderToFile(reader, utility.BackupDBLocation+string(filepath.Separator)+containerName+".sql")

	}
	tarFileLocation := utility.BackupLocation + string(filepath.Separator) + timestampStr + ".tar.gz"
	file.FoldersToTarGZ([]string{utility.BackupDBLocation, utility.ExtraLocation}, tarFileLocation)

	res, err := aws.UploadFile(tarFileLocation, timestampStr)
	panicIfErr(err)

	archiveIDs, err := file.GetArchiveIDs(utility.ArchiveIDLocation)
	panicIfErr(err)

	ttl, err := time.ParseDuration(os.Getenv(utility.ENVTTL))
	panicIfErr(err)

	expires := timestamp.Add(ttl)
	archiveIDs = append(archiveIDs, file.SerializeableArchive{ArchiveID: *res.ArchiveId, Location: *res.Location, Checksum: *res.Checksum, Timestamp: timestamp, Expires: expires})

	err = file.WriteArchiveIDs(archiveIDs, utility.ArchiveIDLocation)
	panicIfErr(err)
}