package configutil

import (
	"context"
	"io/ioutil"
	"strings"
)

// File reads the string contents of a file as a literal config value.
type File string

// String returns the string contents of a file.
func (f File) String(ctx context.Context) (*string, error) {
	contents, err := ioutil.ReadFile(string(f))
	if err != nil {
		return nil, nil
	}
	stringContents := strings.TrimSpace(string(contents))
	return &stringContents, nil
}
