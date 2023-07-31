package monitor

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type DirectorySizer interface {
	GetCurrentSize() (int64, error)
	RemoveElderFiles() error
	MaxSizeBytes() int64
}

type directorySizerImpl struct {
	maxSizeBytes int64
	dirPath      string
}

// NewDirectorySizer creates a new DirectorySizer instance.
//
// It takes in the directory path as a string (`dirPath`) and the maximum size in bytes as an int64 (`MaxSizeBytes`).
// It returns a pointer to a DirectorySizer struct.
func NewDirectorySizer(dirPath string, MaxSizeBytes int64) DirectorySizer {
	return &directorySizerImpl{
		maxSizeBytes: MaxSizeBytes,
		dirPath:      dirPath,
	}
}

// GetCurrentSize returns the current size of the directory.
//
// It takes no parameters.
// It returns an int64, which represents the total size of the directory, and an error if any error occurred during the process.
func (sizer *directorySizerImpl) GetCurrentSize() (int64, error) {
	var totalSize int64

	err := filepath.Walk(sizer.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

// RemoveElderFiles is a function that removes the elder files from the directory.
//
// It reads the directory entries from the specified directory path and sorts them
// by modification time in ascending order. Then, it calculates the total size of
// the files and removes the files from the directory until the total size exceeds
// the maximum size specified. The function returns an error if there is any issue
// reading or removing the files.
//
// Parameters:
// - None
//
// Returns:
// - error: An error if there is any issue reading or removing the files.
func (sizer *directorySizerImpl) RemoveElderFiles() (err error) {
	dirEntries, err := os.ReadDir(sizer.dirPath)
	if err != nil {
		return err
	}

	// Sort files by modification time (oldest first)
	sort.Slice(dirEntries, func(i, j int) bool {
		entryI, _ := dirEntries[i].Info()
		entryJ, _ := dirEntries[j].Info()
		return entryI.ModTime().Before(entryJ.ModTime())
	})

	var totalSize int64
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			inf, _ := entry.Info()
			totalSize += inf.Size()

			if totalSize > sizer.MaxSizeBytes() {
				filePath := filepath.Join(sizer.dirPath, entry.Name())
				if err := os.Remove(filePath); err != nil {
					return err
				}

				//log.Debug().Msgf("removed file: %s", filePath)
			}
		}
	}

	return nil
}

func (sizer *directorySizerImpl) MaxSizeBytes() int64 {
	return sizer.maxSizeBytes
}

type DirectoryMonitor struct {
	sizer    DirectorySizer
	interval time.Duration
}

func NewDirectoryMonitor(dirPath string, MaxSizeBytes int64, interval time.Duration) *DirectoryMonitor {
	sizer := NewDirectorySizer(dirPath, MaxSizeBytes)
	return &DirectoryMonitor{
		sizer:    sizer,
		interval: interval,
	}
}

func (m *DirectoryMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Exit the method when the context is canceled.
		case <-ticker.C:
			currentSize, err := m.sizer.GetCurrentSize()
			if err != nil {
				continue
			}
			if currentSize > m.sizer.MaxSizeBytes() {
				// If the directory size exceeds the threshold.
				m.sizer.RemoveElderFiles()
			}
			// Perform other monitoring tasks as needed.
		}
	}
}
