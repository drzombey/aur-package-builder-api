package storage

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
)

func createTarFileFromStream(filename, dest string, stream io.Reader) (*os.File, error) {
	reader := tar.NewReader(stream)

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		err := os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	filePathWithName := dest + "/" + filename

	file, err := os.Create(filePathWithName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return file, nil
}
