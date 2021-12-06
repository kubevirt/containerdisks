package build

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"time"
)

func StreamLayer(imagePath string) (tarReader io.ReadCloser, errChan chan error) {
	errChan = make(chan error, 1)
	reader, writer := io.Pipe()
	tarWriter := tar.NewWriter(writer)

	go func() {
		defer writer.Close()
		defer tarWriter.Close()
		err := addFileToTarWriter(imagePath, tarWriter)
		if err != nil {
			errChan <- fmt.Errorf("error adding file '%s', to tarball: %v", imagePath, err)
		}
		close(errChan)
	}()

	return reader, errChan
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file information with stat: %v", err)
	}

	header := &tar.Header{
		Typeflag: tar.TypeDir,
		Name:     "disk/",
		Mode:     0555,
		Uid:      107,
		Gid:      107,
		Uname:    "qemu",
		Gname:    "qemu",
		ModTime:  time.Now(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("error writing disks directory tar header: %v", err)
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
		return fmt.Errorf("error writing image file tar header: %v", err)
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("error writingfile into tarball: %v", err)
	}

	return nil
}
