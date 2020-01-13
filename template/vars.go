package template

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/blend/go-sdk/yaml"
)

// Vars is a soft alias to map[string]interface{}.
type Vars = map[string]interface{}

// NewVarsFromPath returns a new vars file from a given path.
func NewVarsFromPath(path string) (map[string]interface{}, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	output := map[string]interface{}{}
	if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(contents, &output)
		if err != nil {
			return nil, err
		}
	} else {
		err = yaml.Unmarshal(contents, &output)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

// MergeVars merges a variadic array of variable sets.
func MergeVars(vars ...Vars) Vars {
	output := Vars{}
	for _, set := range vars {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}
