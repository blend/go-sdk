package fileutil

import (
	"os"
	"time"

	"github.com/blend/go-sdk/ex"
)

// Error constants
const (
	ErrWatchStopped ex.Class = "watch file should stop"
)

// Watch watches a file for changes and calls the action if there are changes.
// It does this by polling the file for ModTime changes every 500ms.
// It is not designed for watching a large number of files.
// You should probably call this within a go routine.
func Watch(path string, action func(*os.File) error) error {
	stat, err := os.Stat(path)
	if err != nil {
		return ex.New(err)
	}

	lastMod := stat.ModTime()

	ticker := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-ticker:
			stat, err = os.Stat(path)
			if err != nil {
				return ex.New(err)
			}
			if stat.ModTime().After(lastMod) {
				file, err := os.Open(path)
				if err != nil {
					return ex.New(err)
				}
				if err := action(file); err != nil {
					if ex.Is(err, ErrWatchStopped) {
						return nil
					}
					return ex.New(err)
				}
				lastMod = stat.ModTime()
			}
		}
	}
}
