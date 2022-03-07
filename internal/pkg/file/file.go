package file

import (
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"path/filepath"
)

// ReaderToFile reads from a docker reader and writes the contents to a file until EOF
func ReaderToFile(reader io.Reader, fileName string) (int64, error) {

	dir := filepath.Dir(fileName)

	os.MkdirAll(dir, 0644)

	// open output files
	filenameErr := fileName + ".err"
	fo, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	foErr, err := os.Create(filenameErr)
	if err != nil {
		return 0, err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	defer func() {
		foErr.Close()
	}()

	n, err := stdcopy.StdCopy(fo, foErr, reader)
	if err != nil {
		return 0, err
	}
	errStat, err := foErr.Stat()
	if err != nil {
		return 0, err
	}
	if errStat.Size() == 0 {
		foErr.Close()
		os.Remove(filenameErr)
	}
	return n, nil
}
