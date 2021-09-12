package file

import (
	"io"
	"os"
	"path/filepath"
)

// ReaderToFile reads from a reader and writes the contents to a file until EOF
func ReaderToFile(reader io.Reader, fileName string) (int, error) {

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
	buf := make([]byte, 1024)
	n_sum := 0
	for {
		// read a chunk
		n, err := reader.Read(buf)
		n_sum += n
		if err != nil && err != io.EOF {
			return n_sum, err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			return n_sum, err
		}
	}
	return n_sum, nil
}
