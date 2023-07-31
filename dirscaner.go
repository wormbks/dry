package dry

import (
	"os"
	"path/filepath"
)

// GetFilesByPattern returns a list of files that match an include pattern and do not match an exclude pattern
// in a given root directory and its subdirectories.
//
// rootPath: the root directory to start searching in.
// includePattern: the pattern that files should match in order to be included in the result list.
// excludePattern: the pattern that files should not match in order to be included in the result list.
//
// returns: a slice of file paths that match the include and exclude patterns.
//
//	An error is returned if there was a problem walking the directory tree.
func GetFilesByPattern(rootPath string, includePattern string, excludePattern string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && matchedByPattern(path, includePattern) && !matchedByPattern(path, excludePattern) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// matchedByPattern returns true if a file name matches a given pattern and false if it doesn't.
// It uses filepath.Match to match the file name against the pattern.
//
// file: the file name to match against the pattern.
// pattern: the pattern to match the file name against.
//
// returns: true if the file name matches the pattern, false otherwise.
func matchedByPattern(file string, pattern string) bool {
	matched, err := filepath.Match(pattern, filepath.Base(file))
	if err != nil {
		return false
	}
	return matched
}
