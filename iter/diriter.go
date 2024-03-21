package iter

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirIterator represents a directory iterator.
type DirIterator struct {
	root      string   // The root directory of the iterator.
	dirs      []string // The list of directories to iterate over.
	current   string   // The current directory in the iteration.
	index     int      // The current index of the iterator.
	indexHash uint64   // The hash value of the current index.
}

// NewDirIterator creates a new directory iterator with the specified root directory.
func NewDirIterator(root string) (*DirIterator, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	dirs, err := getSubdirectories(root)
	if err != nil {
		return nil, err
	}

	return &DirIterator{
		root:      root,
		dirs:      dirs,
		current:   "",
		index:     -1,
		indexHash: calculateIndexHash(dirs),
	}, nil
}

// Next returns the next subdirectory.
func (iter *DirIterator) Next() (string, error) {
	dirs, err := getSubdirectories(iter.root)
	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		return "", fmt.Errorf("no subdirectories found in %s", iter.root)
	}

	// Check if the list of subdirectories has changed
	if iter.indexHash != calculateIndexHash(dirs) {
		// Update the list of subdirectories
		iter.dirs = dirs
		iter.indexHash = calculateIndexHash(dirs)
	}

	// Move to the next index, loop back if necessary
	iter.index = (iter.index + 1) % len(iter.dirs)

	// Get the current subdirectory
	iter.current = iter.dirs[iter.index]

	return iter.current, nil
}

// getSubdirectories returns a list of subdirectories.
func getSubdirectories(root string) ([]string, error) {
	var dirs []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			dirs = append(dirs, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort the directories for consistent ordering
	sort.Strings(dirs)

	return dirs, nil
}

// calculateIndexHash calculates the FNV-1a hash for the list of subdirectories.
func calculateIndexHash(dirs []string) uint64 {
	// Concatenate the sorted directory names
	concatenated := strings.Join(dirs, ",")

	// Calculate the FNV-1a hash
	hasher := fnv.New64a()
	hasher.Write([]byte(concatenated)) // #nosec G104 It returns nil error all the time.
	hashValue := hasher.Sum64()

	return hashValue
}
