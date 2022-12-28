package storage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"os"
	"testing"
)

var testFolder = os.TempDir() + "/test_files"

func SetupTest() {
	err := os.MkdirAll(testFolder, os.ModePerm)
	if err != nil {
		return
	}
}

func createTestTarGzFile(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	hdr := &tar.Header{
		Name: "test.tar",
		Mode: 0644,
		Size: int64(len(content)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	_, err = tw.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func TestFileSystemProvider_GetFile(t *testing.T) {
	SetupTest()
	provider := NewFilesystemProvider(testFolder)
	defer os.RemoveAll(testFolder)

	ctx := context.Background()
	fileName := "test.txt"
	content := []byte("This is a test file")
	err := os.WriteFile(testFolder+"/"+fileName, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file, err := provider.GetFile(ctx, fileName)
	if err != nil {
		t.Fatalf("Failed to get file: %v", err)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(fileData, content) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", content, fileData)
	}
}

func TestFileSystemProvider_DeleteFile(t *testing.T) {
	SetupTest()
	provider := NewFilesystemProvider(testFolder)
	defer os.RemoveAll(testFolder)

	ctx := context.Background()
	fileName := "test.txt"
	content := []byte("This is a test file")
	err := os.WriteFile(testFolder+"/"+fileName, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = provider.DeleteFile(ctx, fileName)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	_, err = os.Stat(testFolder + "/" + fileName)
	if !os.IsNotExist(err) {
		t.Errorf("File was not deleted: %v", err)
	}
}
