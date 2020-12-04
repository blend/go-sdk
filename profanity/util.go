package profanity

import "fmt"

func makeStrings(objs ...interface{}) []string {
	var output []string
	for _, obj := range objs {
		output = append(output, fmt.Sprint(obj))
	}
	return output
}
