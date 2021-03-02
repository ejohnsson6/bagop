package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"
)

// SerializeableArchive is an archive which can be serialized as a CSV object
type SerializeableArchive struct {
	ArchiveID        string    `csv:"archive_id"`
	Checksum         string    `csv:"checksum"`
	Location         string    `csv:"location"`
	Timestamp        time.Time `csv:"timestamp"`
	Expires          bool      `csv:"expires"`
	ExpiresTimestamp time.Time `csv:"expires_timestamp"`
}

func GetArchiveIDs(csvFile string) ([]SerializeableArchive, error) {

	if err := os.MkdirAll(filepath.Dir(csvFile), 0644); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(csvFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	var out []SerializeableArchive
	if err = gocsv.UnmarshalFile(f, &out); err != nil && err != gocsv.ErrEmptyCSVFile {
		return nil, err
	}

	return out, nil

}

func WriteArchiveIDs(archiveIDs []SerializeableArchive, csvFile string) error {

	if err := os.MkdirAll(filepath.Base(csvFile), 0644); err != nil {
		return err
	}
	f, err := os.OpenFile(csvFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	gocsv.MarshalFile(&archiveIDs, f)

	return nil

}
