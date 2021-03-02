package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"
)

// SerializeableArchive is an archive which can be serialized as a CSV object
type SerializeableArchive struct {
	ArchiveID string    `csv:"archive_id"`
	Timestamp time.Time `csv:"timestamp"`
	Expires   time.Time `csv:"expires"`
}

func GetArchiveIDs(csvFile string) []SerializeableArchive {

	return nil

}

func WriteArchiveIDs(archiveIDs []SerializeableArchive, csvFile string) error {

	b, err := gocsv.MarshalBytes(archiveIDs)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Base(csvFile), 0644); err != nil {
		return err
	}
	f, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(b)

	return nil

}
