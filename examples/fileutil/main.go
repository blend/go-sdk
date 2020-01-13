package main

import (
	"fmt"
	"os"
	"github.com/blend/go-sdk/fileutil"
)

func main() {
	go fileutil.Watch("file.yml", func(f *os.File) error {
		defer f.Close()
		fmt.Printf("file changed\n")
		return nil
	})

	select {}
}
