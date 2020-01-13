package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blend/go-sdk/sh"
	"github.com/blend/go-sdk/yaml"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type configVar struct {
	Field   string
	Default string
	Value   string
}

func (cv *configVar) Set(value string) error {
	cv.Value = value
	return nil
}

type configVarSet []configVar

func (cvs configVarSet) String() string {
	return "Fields to prompt for in the form fieldname=default"
}

func (cvs configVarSet) Union(other configVarSet) configVarSet {
	return append(cvs, other...)
}

func (cvs *configVarSet) Set(flagValue string) error {
	parts := strings.SplitN(flagValue, "=", 2)
	if len(parts) > 1 {
		*cvs = append(*cvs, configVar{
			Field:   parts[0],
			Default: parts[1],
		})
	} else if len(parts) > 0 {
		*cvs = append(*cvs, configVar{
			Field: parts[0],
		})
	} else {
		return fmt.Errorf("invalid config var: %s", flagValue)
	}
	return nil
}

func (cvs configVarSet) MarshalYAML() (interface{}, error) {
	output := make(map[string]interface{})
	for _, cv := range cvs {
		if cv.Value != "" {
			output[cv.Field] = cv.Value
		} else {
			output[cv.Field] = cv.Default
		}
	}
	return output, nil
}

func main() {
	fields := configVarSet{}
	secureFields := configVarSet{}
	output := flag.String("o", "", "The output file")
	flag.Var(&fields, "field", "Fields to prompt for, can be multiple")
	flag.Var(&secureFields, "secure", "Secure Fields to prompt for, can be multiple, will hide input")
	flag.Parse()

	if len(fields) == 0 && len(secureFields) == 0 {
		fmt.Fprintf(os.Stderr, "please provide at least (1) field or a secure field\n")
		flag.Usage()
	}

	for index := range fields {
		prompt(&fields[index])
	}

	for index := range secureFields {
		secure(&secureFields[index])
	}

	all := fields.Union(secureFields)
	buffer := new(bytes.Buffer)
	if err := yaml.NewEncoder(buffer).Encode(all); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := ioutil.WriteFile(*output, buffer.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprint(os.Stderr, buffer.String())
	}
}

func prompt(cv *configVar) {
	if len(cv.Default) > 0 {
		cv.Value = sh.Promptf("%s[%s]: ", cv.Field, cv.Default)
	} else {
		cv.Value = sh.Promptf("%s: ", cv.Field)
	}
}

func secure(cv *configVar) {
	var prompt string
	if len(cv.Default) > 0 {
		prompt = fmt.Sprintf("%s[%s]: ", cv.Field, cv.Default)
	} else {
		prompt = fmt.Sprintf("%s: ", cv.Field)
	}
	output := sh.MustPassword(prompt)
	cv.Value = output
}
