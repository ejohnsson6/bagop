package file

import (
	"os"
	"path/filepath"
)

// ReaderToFile reads from a reader and writes the contents to a file until EOF
func StringToFile(str string, fileName string) (int, error) {

	dir := filepath.Dir(fileName)

	os.MkdirAll(dir, 0644)

	// open output file
	fo, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	n, err := fo.WriteString(str)
	return n, nil
}
