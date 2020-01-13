package fileutil

import (
	"bufio"
	"os"
	"github.com/blend/go-sdk/ex"
)

// ReadLines reads a file and calls the handler for each line.
func ReadLines(filePath string, handler func(string) error) error {
	f, err := os.Open(filePath)
	if err != nil {
		return ex.New(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		err = handler(line)
		if err != nil {
			return ex.New(err)
		}
	}
	return nil
}
