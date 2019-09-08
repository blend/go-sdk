package bindata

import "sort"

// Assets represent a path => File
type Assets map[string]*File

// Paths returns the file paths sorted by name.
func (a Assets) Paths() []string {
	var paths []string
	for key := range a {
		paths = append(paths, key)
	}
	sort.Strings(paths)
	return paths
}
