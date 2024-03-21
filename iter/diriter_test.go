package iter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirIteratorWithTempDir(t *testing.T) {
	// Create a temporary test directory
	rootDir, err := os.MkdirTemp("", "prefix")
	assert.NoError(t, err, "Error creating temporary test directory")
	defer os.RemoveAll(rootDir)

	// Create subdirectories in the temporary test directory
	subDir1 := filepath.Join(rootDir, "subdir1")
	subDir2 := filepath.Join(rootDir, "subdir2")
	err = os.Mkdir(subDir1, 0o755)
	assert.NoError(t, err, "Error creating subdirectory 1")
	err = os.Mkdir(subDir2, 0o755)
	assert.NoError(t, err, "Error creating subdirectory 2")

	// Create DirIterator with the temporary test directory
	iter, err := NewDirIterator(rootDir)
	assert.NoError(t, err, "Error creating DirIterator")
	assert.NotNil(t, iter, "DirIterator should not be nil")

	// Test Next method
	for i := 0; i < 2; i++ {
		dir, err := iter.Next()
		assert.NoError(t, err, "Error calling Next()")
		assert.NotEmpty(t, dir, "Directory should not be empty")
	}

	// Test case where the root directory changes (simulate directory changes during iteration)
	newSubDir := filepath.Join(rootDir, "newsubdir")
	err = os.Mkdir(newSubDir, 0o755)
	assert.NoError(t, err, "Error creating new subdirectory")

	// Test Next method after directory changes
	dir, err := iter.Next()
	assert.NoError(t, err, "Error calling Next() after directory changes")
	assert.NotEmpty(t, dir, "Directory should not be empty after changes")

	// Test case where the root directory is removed
	err = os.RemoveAll(rootDir)
	assert.NoError(t, err, "Error removing root directory")

	// Test Next method after root directory removal
	dir, err = iter.Next()
	assert.Error(t, err, "Error calling Next() after root directory removal")
	assert.Empty(t, dir, "Directory should be empty after root directory removal")
}
