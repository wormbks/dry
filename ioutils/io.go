package ioutils

import (
	"fmt"
	"os"
)

// CreateFolderStructure creates a folder structure based on the provided path.
// It creates all necessary parent directories if they don't exist.
func CreateFolderStructure(path string) error {
	// os.ModePerm ensures that the folders are created with read, write,
	// and execute permissions for the current user.
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder structure: %w", err)
	}
	return nil
}
