/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/blend/go-sdk/template"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
)

// Includes are a collection of template files to include as sub templates.
type Includes []string

// Set sets the value.
func (v *Includes) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *Includes) String() string {
	return "Files to include as sub templates"
}

// Variables are a list of commandline variables.
type Variables []string

// Set sets a variable.
func (v *Variables) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *Variables) String() string {
	return "Variable values to set in the template"
}

// Values returns the map of values.
func (v *Variables) Values() (values map[string]string) {
	values = map[string]string{}

	for _, val := range *v {
		pieces := strings.SplitN(val, "=", 2)
		if len(pieces) > 1 {
			values[pieces[0]] = pieces[1]
		}
	}
	return
}

func main() {
	var templateFile string
	flag.StringVar(&templateFile, "f", "", "Template file to process; if \"-\", will read from os.Stdin")

	var includes Includes
	flag.Var(&includes, "i", "Files to include as sub templates")

	var varsFile string
	flag.StringVar(&varsFile, "vars", "", "Vars file to process")

	var outFile string
	flag.StringVar(&outFile, "o", "", "Output file")

	var delims string
	flag.StringVar(&delims, "delims", "", "Delimiters in the form --delims=\"{{,}}\"")

	var variables Variables
	flag.Var(&variables, "var", "Variables in the form --var=foo=bar")

	var help bool
	flag.BoolVar(&help, "help", false, "Shows this usage message")

	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Shows the app version")

	flag.Usage = func() {
		if len(version) == 0 {
			version = "master"
		}
		fmt.Fprintf(os.Stderr, "%s version %s\n\n", os.Args[0], version)
		fmt.Fprintf(os.Stderr, "Find more information at https://github.com/blend/go-sdk/tree/master/template\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample Usage:\n")
		fmt.Fprintf(os.Stderr, "Read a template file: \"template -f template.yml\"\n")
		fmt.Fprintf(os.Stderr, "Read a template from stdin: \"echo '{{ .Var \"foo\" }}' | template -f -\"\n")
		fmt.Fprintf(os.Stderr, "Specify a variable: template -f config.yml --var=foo=bar\n")
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		if len(version) == 0 {
			version = "master"
		}
		fmt.Fprintf(os.Stdout, "%s version %s %s/%s\n", os.Args[0], version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	var temp *template.Template
	var err error
	if len(templateFile) > 0 && templateFile == "-" {
		temp = template.New()

		var contents []byte
		contents, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		temp = temp.WithBody(string(contents))
	} else if len(templateFile) > 0 {
		temp, err = template.NewFromFile(templateFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}

	if len(includes) > 0 {
		for _, include := range includes {
			var contents []byte
			contents, err = os.ReadFile(include)
			if err != nil {
				log.Fatal(err)
			}
			temp = temp.WithInclude(string(contents))
		}
	}

	if len(varsFile) > 0 {
		vars, err := template.NewVarsFromPath(varsFile)
		if err != nil {
			log.Fatal(err)
		}
		for key, value := range vars {
			temp = temp.WithVar(key, value)
		}
	}

	vars := variables.Values()
	if len(vars) > 0 {
		for key, value := range vars {
			temp = temp.WithVar(key, value)
		}
	}

	if len(delims) > 0 {
		d := strings.Split(delims, ",")
		if len(d) < 2 {
			log.Fatalf("Invalid delimiters: %s", delims)
		}
		temp.WithDelims(d[0], d[1])
	}

	buffer := bytes.NewBuffer(nil)
	err = temp.Process(buffer)
	if err != nil {
		log.Fatal(err)
	}

	if len(outFile) > 0 {
		f, err := os.Create(outFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = buffer.WriteTo(f)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, _ = buffer.WriteTo(os.Stdout)
	}
}
