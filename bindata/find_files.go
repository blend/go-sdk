package bindata

import (
	"os"
	"path/filepath"
	"regexp"
)

// FindFiles traverses a root recursively, ignoring files that match the optional ignore
// expression, and calls the handler when it finds a file.
func FindFiles(root string, ignores []*regexp.Regexp, handler func(*File)) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, ignore := range ignores {
			if ignore.MatchString(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		f, err := ReadFile(path)
		if err != nil {
			return err
		}
		handler(f)
		return nil
	})
}
