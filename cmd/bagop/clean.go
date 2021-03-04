package main

import (
	"time"

	"github.com/swexbe/bagop/internal/pkg/aws"
	"github.com/swexbe/bagop/internal/pkg/file"
	l "github.com/swexbe/bagop/internal/pkg/logging"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

// Returns true if archive has expired
func hasExpired(archive file.SerializeableArchive) bool {
	return archive.Expires && archive.ExpiresTimestamp.Before(time.Now())
}

// Tries to delete the archive from the vault if expired
// Returns true if archive should be kept in storage
// i.e. if deletion failed or not expired
func filterArchiveHelper(archive file.SerializeableArchive, vaultName string) bool {
	l.Logger.Debugf("Archive %s expires: %t (%s)", archive.ArchiveID, archive.Expires, archive.ExpiresTimestamp.Format(time.RFC3339))
	if hasExpired(archive) {
		l.Logger.Infof("Archive %s has expired, deleting from Glacier...", archive.ArchiveID)
		_, err := aws.DeleteArchive(vaultName, archive.ArchiveID)
		if err != nil {
			l.Logger.Warnf("Error when deleting archive: %s", err.Error())
			return true
		}
		l.Logger.Debugf("Deletion succeeded")
		return false
	}
	return true
}

func cleanBackups(vaultName string) {

	l.Logger.Infof("Checking for expired archives")
	archives, err := file.GetArchiveIDs(utility.ArchiveIDLocation)
	panicIfErr(err)
	l.Logger.Debugf("Going through %d archives", len(archives))

	numDeleted := 0

	var archivesNew []file.SerializeableArchive

	for _, archive := range archives {
		if filterArchiveHelper(archive, vaultName) {
			archivesNew = append(archivesNew, archive)
		} else {
			numDeleted++
		}
	}

	file.WriteArchiveIDs(archivesNew, utility.ArchiveIDLocation)
	panicIfErr(err)
	l.Logger.Infof("Finished, %d archives deleted", numDeleted)
}
