package file

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteStringToFile appends a new line and string to a file
// Creating the file and directory if it doesn't exist
func WriteStringToFile(fileName string, perm fs.FileMode, text string) error {

	err := os.MkdirAll(filepath.Dir(fileName), perm)

	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString("\n" + text); err != nil {
		return err
	}
	return nil
}

// ReaderToFile reads from a reader and writes the contents to a file until EOF
func ReaderToFile(reader io.Reader, fileName string) error {

	dir := filepath.Dir(fileName)

	os.MkdirAll(dir, 0755)

	// open output file
	fo, err := os.Create(fileName)
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
