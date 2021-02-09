package file

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	l "github.com/swexbe/bagop/internal/pkg/logging"
)

// FolderToTarGZ compresses all files in a directory into a gzipped tar file
func FolderToTarGZ(sourcedir string, destinationfile string) error {

	dir, err := os.Open(sourcedir)

	if err != nil {
		return err
	}

	defer dir.Close()

	files, err := dir.Readdir(0) // grab the files list

	if err != nil {
		return err
	}

	tarfile, err := os.Create(destinationfile)
	l.Logger.Infof("Created Tarball %s", tarfile.Name())

	if err != nil {
		return err
	}

	defer tarfile.Close()
	var fileWriter io.WriteCloser = tarfile

	fileWriter = gzip.NewWriter(tarfile)
	defer fileWriter.Close()
	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {
		l.Logger.Infof("Compressing file %s", fileInfo.Name())

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())

		if err != nil {
			return err
		}

		defer file.Close()

		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)

		if err != nil {
			return err
		}

		_, err = io.Copy(tarfileWriter, file)

		if err != nil {
			return err
		}
	}
	return nil
}
