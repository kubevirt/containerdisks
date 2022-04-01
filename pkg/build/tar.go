package build

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"time"
)

func StreamLayerOpener(imagePath string) func() (io.ReadCloser, error) {
	modTime := time.Now()

	return func() (io.ReadCloser, error) {
		fileErrorChan := make(chan error)
		pipeReader, pipeWriter := io.Pipe()

		go func() {
			defer pipeWriter.Close()

			file, err := os.Open(imagePath)
			if err != nil {
				fileErrorChan <- fmt.Errorf("error opening file: %w", err)
				return
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				fileErrorChan <- fmt.Errorf("error getting file information with stat: %w", err)
				return
			}

			// Close channel after successfully opening file to avoid deadlock
			close(fileErrorChan)

			tarWriter := tar.NewWriter(pipeWriter)
			err = addFileToTarWriter(file, stat, modTime, tarWriter)
			if err != nil {
				// Move the error to the PipeReader side. It is ok to call close on PipeWriter multiple times.
				pipeWriter.CloseWithError(fmt.Errorf("error adding file '%s', to tarball: %w", imagePath, err))
			}
			err = tarWriter.Close()
			if err != nil {
				pipeWriter.CloseWithError(fmt.Errorf("error writing footer of tarball: %w", err))
			}
		}()

		// Wait until file is opened or immediately return any errors
		if err, ok := <-fileErrorChan; ok {
			return nil, err
		}

		return pipeReader, nil
	}
}

func addFileToTarWriter(file io.Reader, stat os.FileInfo, modTime time.Time, tarWriter *tar.Writer) error {
	header := &tar.Header{
		Typeflag: tar.TypeDir,
		Name:     "disk/",
		Mode:     0555,
		Uid:      107,
		Gid:      107,
		Uname:    "qemu",
		Gname:    "qemu",
		ModTime:  modTime,
	}

	err := tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("error writing disks directory tar header: %w", err)
	}

	header = &tar.Header{
		Typeflag: tar.TypeReg,
		Uid:      107,
		Gid:      107,
		Uname:    "qemu",
		Gname:    "qemu",
		Name:     "disk/disk.img",
		Size:     stat.Size(),
		Mode:     0444,
		ModTime:  stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("error writing image file tar header: %w", err)
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("error writingfile into tarball: %w", err)
	}

	return nil
}
