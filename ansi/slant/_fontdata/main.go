/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	contents, err := os.ReadFile("slant_letters.txt")
	if err != nil {
		log.Fatal(err)
	}

	letterHeight := 6

	scanner := bufio.NewScanner(bytes.NewReader(contents))

	var lineText string
	var line int

	fmt.Fprintln(os.Stdout, "[][]string{")

	for scanner.Scan() {
		lineText = scanner.Text()

		if line == 0 {
			fmt.Fprintln(os.Stdout, "\t{")
		}
		if line < letterHeight {
			lineText = strings.TrimSuffix(strings.TrimSuffix(lineText, "@"), "@")
			lineText = strings.ReplaceAll(lineText, "\"", "\\\"")	// escape quotes
			lineText = strings.ReplaceAll(lineText, "\\", "\\\\")	// escape slashes
			fmt.Fprintf(os.Stdout, "\t\t\"%s\",\n", lineText)
			line++
		}

		if line == letterHeight || strings.HasSuffix(lineText, "@@") {
			fmt.Fprintln(os.Stdout, "\t},")
			line = 0
		}
	}

	fmt.Fprintln(os.Stdout, "}")
}
