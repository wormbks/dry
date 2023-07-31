package async

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AsyncGzipFileWriter(t *testing.T) {
	// Create a temporary file for testing.
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create an AsyncGzipFileWriter for testing.
	writer, err := NewAsyncGzipFileWriter(context.Background(), tempFile.Name(), true)
	if err != nil {
		t.Fatalf("Failed to create AsyncGzipFileWriter: %v", err)
	}
	defer writer.Close()
	wg := &sync.WaitGroup{}
	writer.Start(wg)

	// Write data to the writer.
	data := []byte("Hello, World!")
	n, err := writer.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// Close the writer.
	err = writer.Close()
	assert.NoError(t, err)
	wg.Wait()

	// Verify that the file contains the compressed data.
	fileContent, err := ioutil.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Decompress the file content for verification.
	gzipReader, err := gzip.NewReader(bytes.NewReader(fileContent))
	assert.NoError(t, err)
	defer gzipReader.Close()

	decompressedData, err := ioutil.ReadAll(gzipReader)
	assert.NoError(t, err)

	// Verify that the decompressed data matches the original data.
	assert.Equal(t, data, decompressedData)
}

func Test_AsyncGzipFileWriter_ErrorCreateFile(t *testing.T) {
	// Try to create an AsyncGzipFileWriter with a non-existent directory.
	// This should trigger an error during file creation.
	_, err := NewAsyncGzipFileWriter(context.Background(), "/nonexistent_directory/test_file", true)
	assert.Error(t, err, "Expected error when creating file in non-existent directory")

	// Verify that the error is due to non-existent directory.
	assert.True(t, os.IsNotExist(err), "Expected error due to non-existent directory")
}

func Test_AsyncGzipFileWriter_NoGzip(t *testing.T) {
	// Create a temporary file for testing.
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create an AsyncGzipFileWriter for testing.
	writer, err := NewAsyncGzipFileWriter(context.Background(), tempFile.Name(), true)
	if err != nil {
		t.Fatalf("Failed to create AsyncGzipFileWriter: %v", err)
	}
	defer writer.Close()
	wg := &sync.WaitGroup{}
	writer.Start(wg)

	// Write data to the writer.
	data := []byte("Hello, World!")
	n, err := writer.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// Close the writer.
	err = writer.Close()
	assert.NoError(t, err)
	wg.Wait()

	// Verify that the file contains the compressed data.
	fileContent, err := ioutil.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	// Decompress the file content for verification.
	gzipReader, err := gzip.NewReader(bytes.NewReader(fileContent))
	assert.NoError(t, err)
	defer gzipReader.Close()

	decompressedData, err := ioutil.ReadAll(gzipReader)
	assert.NoError(t, err)

	// Verify that the decompressed data matches the original data.
	assert.Equal(t, data, decompressedData)
}
