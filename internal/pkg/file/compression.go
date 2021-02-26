package file

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	l "github.com/swexbe/bagop/internal/pkg/logging"
)

func recursiveToTarGZ(sourcedir string, tarfileWriter *tar.Writer) error {

	dir, err := os.Open(sourcedir)

	if err != nil {
		return err
	}

	defer dir.Close()

	files, err := dir.Readdir(0) // grab the files list

	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		l.Logger.Infof("Compressing file %s", fileInfo.Name())

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if fileInfo.IsDir() {
			l.Logger.Infof("Recursively compressing files in folder %s", file.Name())
			recursiveToTarGZ(file.Name(), tarfileWriter)
			continue
		}

		if err != nil {
			return err
		}

		defer file.Close()

		header := new(tar.Header)
		header.Name = file.Name()
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

// FoldersToTarGZ compresses all files in a list of directores into a gzipped tar file
// Paths are kept the same as the full system path to avoid conflicts
func FoldersToTarGZ(sourcedirs []string, destinationfile string) error {

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

	for _, sourcedir := range sourcedirs {
		l.Logger.Infof("Compressing files in folder %s", sourcedir)
		err = recursiveToTarGZ(sourcedir, tarfileWriter)

		if err != nil {
			return err
		}

	}

	return nil
}
