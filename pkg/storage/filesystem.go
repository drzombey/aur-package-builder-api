package storage

import (
	"context"
	"io"
	"io/fs"
	"os"
)

type FileSystemProvider struct {
	dest string
}

func NewFilesystemProvider(dest string) Provider {
	return &FileSystemProvider{dest: dest}
}

func (f FileSystemProvider) AddFile(ctx context.Context, fileName string, stream io.Reader) error {
	_, err := io.ReadAll(stream)
	if err != nil {
		return err
	}

	_, err = createTarFileFromStream(fileName, f.dest, stream)
	if err != nil {
		return err
	}
	return nil
}

func (f FileSystemProvider) GetFile(ctx context.Context, filename string) (*os.File, error) {
	return os.OpenFile(f.dest+"/"+filename, os.O_RDONLY, fs.ModePerm)
}

func (f FileSystemProvider) DeleteFile(ctx context.Context, filename string) error {
	return os.Remove(f.dest + "/" + filename)
}
