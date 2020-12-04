package db

import "encoding/json"

// JSON returns the json representation of a given object for inserts / updates.
func JSON(obj interface{}) interface{} {
	jsonBytes, _ := json.Marshal(obj)
	if result := string(jsonBytes); result != "null" { // explicitly bad.
		return result
	}
	return nil
}
