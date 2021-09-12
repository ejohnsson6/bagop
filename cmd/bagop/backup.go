package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/swexbe/bagop/internal/pkg/aws"
	"github.com/swexbe/bagop/internal/pkg/db"
	"github.com/swexbe/bagop/internal/pkg/docker"
	"github.com/swexbe/bagop/internal/pkg/file"
	l "github.com/swexbe/bagop/internal/pkg/logging"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

func getDumpCmd(vendor string, env []string) []string {
	switch vendor {
	case "mysql", "mariadb":
		l.Logger.Debugf("Dumping as MYSQL/MariaDB container")
		return db.DumpMysqlCmd(env)
	case "postgres":
		l.Logger.Debugf("Dumping as PostgreSQL container")
		return db.DumpPostgresCmd(env)
	}
	// Should never happen
	panic(fmt.Errorf("No dump command for this vendor: %s", vendor))
}

func parseExpirationDate(now time.Time, ttl string) (bool, time.Time) {
	if ttl == "" {
		return false, time.Time{}
	}
	ttlInt, err := strconv.Atoi(ttl)
	// Default to never if TTL can't be parsed
	if err != nil {
		l.Logger.Warnf("Couldn't parse TTL, defaulting to no expiration: %s", err.Error())
		return false, time.Time{}
	}
	expiresTimestamp := now.AddDate(0, 0, ttlInt)
	l.Logger.Debugf("Archive will expire %s", expiresTimestamp.Format(time.RFC3339))
	return true, expiresTimestamp

}

func makeBackup(ttl string, vaultName string) {
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
		l.Logger.Debugf("Trying to dump container %s", container.ID[0:12])
		containerName := docker.FindName(container)
		l.Logger.Infof("Dumping as %s", containerName)
		vendor, err := docker.FindVendor(cli, container)
		if err != nil {
			l.Logger.Errorf(err.Error())
			continue
		}
		env := docker.GetEnv(cli, container)
		cmd := getDumpCmd(vendor, env)

		panicIfErr(err)
		exitCode, str, err := docker.RunCommand(cli, container, cmd)
		panicIfErr(err)
		l.Logger.Debugf("Dump process exited with code: %d", exitCode)
		if exitCode != 0 {
			l.Logger.Errorf("Exit code not 0, run with -v to see output")
			l.Logger.Debug()
			continue
		}
		fileName := utility.BackupDBLocation + string(filepath.Separator) + containerName + ".sql"
		n, err := file.StringToFile(str, fileName)
		panicIfErr(err)
		l.Logger.Debugf("Wrote %d bytes to file: %s", n, fileName)
	}

	tarFileLocation := utility.BackupLocation + string(filepath.Separator) + timestampStr + ".tar.gz"
	file.FoldersToTarGZ([]string{utility.BackupDBLocation, utility.ExtraLocation}, tarFileLocation)
	result, err := aws.UploadFile(tarFileLocation, timestampStr, vaultName)
	panicIfErr(err)

	l.Logger.Debugf("Writing archive id to csv: %s", *result.ArchiveId)
	expires, expiresTimestamp := parseExpirationDate(timestamp, ttl)

	archiveIDs, err := file.GetArchiveIDs(utility.ArchiveIDLocation)
	panicIfErr(err)
	archiveIDs = append(archiveIDs, file.SerializeableArchive{
		ArchiveID:        *result.ArchiveId,
		Location:         *result.Location,
		Checksum:         *result.Checksum,
		Timestamp:        timestamp,
		Expires:          expires,
		ExpiresTimestamp: expiresTimestamp,
	})
	err = file.WriteArchiveIDs(archiveIDs, utility.ArchiveIDLocation)
	panicIfErr(err)

	l.Logger.Infof("Backup succeeded")
}
