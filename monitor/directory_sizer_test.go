package monitor

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDirectorySizer_GetCurrentSize(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := ioutil.TempDir("", "test_directory")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test files in the directory
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(tmpDir, "file3.txt")

	data := []byte("Test data")
	if err := ioutil.WriteFile(file1, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 1: %s", err)
	}
	if err := ioutil.WriteFile(file2, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 2: %s", err)
	}
	if err := ioutil.WriteFile(file3, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 3: %s", err)
	}

	// Initialize the DirectorySizer
	sizer := NewDirectorySizer(tmpDir, 100) // Setting threshold to 100 bytes for testing

	// Test GetCurrentSize
	currentSize, err := sizer.GetCurrentSize()
	assert.Nil(t, err, "GetCurrentSize should not return an error")
	assert.Equal(t, int64(3*len(data)), currentSize, "Incorrect current directory size")
	// Negative test case: Invalid directory path
	invalidDir := "/non_existent_directory"
	invalidSizer := NewDirectorySizer(invalidDir, 100)
	_, err = invalidSizer.GetCurrentSize()
	assert.NotNil(t, err, "GetCurrentSize with invalid directory should return an error")
}

func TestDirectorySizer_RemoveElderFiles(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := ioutil.TempDir("", "test_directory")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test files in the directory
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(tmpDir, "file3.txt")

	data := []byte("Test data")
	if err := ioutil.WriteFile(file1, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 1: %s", err)
	}
	if err := ioutil.WriteFile(file2, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 2: %s", err)
	}
	if err := ioutil.WriteFile(file3, data, 0o644); err != nil {
		t.Fatalf("Error creating test file 3: %s", err)
	}

	// Initialize the DirectorySizer
	sizer := NewDirectorySizer(tmpDir, int64(3*len(data)-1)) // Set threshold to be slightly below the total size of files

	// Test RemoveElderFiles
	err = sizer.RemoveElderFiles()
	assert.Nil(t, err, "RemoveElderFiles should not return an error")

	// After removing the elder files, the directory size should be below the threshold
	currentSize, err := sizer.GetCurrentSize()
	assert.Nil(t, err, "GetCurrentSize should not return an error")
	assert.Less(t, currentSize, int64(3*len(data)), "Directory size should be below the threshold after removing elder files")

	// Negative test case: Directory not found
	invalidDir := "/non_existent_directory"
	invalidSizer := NewDirectorySizer(invalidDir, 100)
	err = invalidSizer.RemoveElderFiles()
	assert.NotNil(t, err, "RemoveElderFiles with invalid directory should return an error")
}

// MockDirectorySizer is a mock implementation of the DirectorySizer interface for testing.
type MockDirectorySizer struct {
	mock.Mock
	// currentSize  int64
	// maxSizeBytes int64
}

func (m *MockDirectorySizer) GetCurrentSize() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDirectorySizer) RemoveElderFiles() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDirectorySizer) MaxSizeBytes() int64 {
	args := m.Called()
	return args.Get(1).(int64)
}

func TestDirectoryMonitor_Start(t *testing.T) {
	// Create a mock DirectorySizer and set the expected return values.
	mockSizer := new(MockDirectorySizer)
	mockSizer.On("GetCurrentSize").Return(int64(100), nil).Times(3)
	mockSizer.On("GetCurrentSize").Return(int64(200), nil).Times(2)
	mockSizer.On("MaxSizeBytes").Return(int64(300)).Times(1)

	// Initialize the DirectoryMonitor with the mock DirectorySizer.
	monitor := &DirectoryMonitor{
		sizer:    mockSizer,
		interval: 1 * time.Second,
	}

	// Create a context with cancel to control the monitoring process.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the monitoring process in a separate goroutine.
	go monitor.Start(ctx)

	// Wait for the monitoring process to run for a few iterations.
	time.Sleep(5 * time.Second)

	// Cancel the context to stop the monitoring process.
	cancel()

	// Assert that the mock methods were called as expected.
	mockSizer.AssertExpectations(t)
}

func TestDirectoryMonitor_Start_Error(t *testing.T) {
	// Create a mock DirectorySizer and set the expected error return value.
	mockSizer := new(MockDirectorySizer)
	mockError := assert.AnError
	mockSizer.On("GetCurrentSize").Return(int64(0), mockError)

	// Initialize the DirectoryMonitor with the mock DirectorySizer.
	monitor := &DirectoryMonitor{
		sizer:    mockSizer,
		interval: 1 * time.Second,
	}

	// Create a context with cancel to control the monitoring process.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the monitoring process in a separate goroutine.
	go monitor.Start(ctx)

	// Wait for the monitoring process to run for a few iterations.
	time.Sleep(3 * time.Second)

	// Cancel the context to stop the monitoring process.
	cancel()

	// Assert that the mock method was called as expected and returned the error.
	mockSizer.AssertExpectations(t)
}

func TestDirectoryMonitor_Start_Cancel(t *testing.T) {
	// Create a mock DirectorySizer with no expectations.
	mockSizer := new(MockDirectorySizer)

	// Initialize the DirectoryMonitor with the mock DirectorySizer.
	monitor := &DirectoryMonitor{
		sizer:    mockSizer,
		interval: 1 * time.Second,
	}

	// Create a context with cancel to control the monitoring process.
	ctx, cancel := context.WithCancel(context.Background())

	// Start the monitoring process in a separate goroutine.
	go monitor.Start(ctx)

	// Wait for a short duration and then cancel the context.
	time.Sleep(2 * time.Second)
	cancel()

	// Assert that the mock methods were not called, as the monitoring should have stopped.
	mockSizer.AssertExpectations(t)
}
