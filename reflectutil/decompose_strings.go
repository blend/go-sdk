/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reflectutil

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
)

// DecomposeStrings decomposes an object into a string map.
func DecomposeStrings(obj interface{}, tagName ...string) map[string]string {
	output := map[string]string{}

	objMeta := reflectType(obj)
	objValue := reflectValue(obj)

	var field reflect.StructField
	var fieldValue reflect.Value
	var tag, tagValue string
	var dataField string
	var pieces []string
	var isCSV bool
	var isBytes bool
	var isBase64 bool

	if len(tagName) > 0 {
		tag = tagName[0]
	}

	for x := 0; x < objMeta.NumField(); x++ {
		isCSV = false
		isBytes = false
		isBase64 = false

		field = objMeta.Field(x)
		if !IsExported(field.Name) {
			continue
		}

		fieldValue = objValue.FieldByName(field.Name)
		dataField = field.Name

		if field.Type.Kind() == reflect.Struct {
			childFields := DecomposeStrings(fieldValue.Interface(), tagName...)
			for key, value := range childFields {
				output[key] = value
			}
		}

		if len(tag) > 0 {
			tagValue = field.Tag.Get(tag)
			if len(tagValue) > 0 {
				if field.Type.Kind() == reflect.Map {
					continue
				} else {
					pieces = strings.Split(tagValue, ",")
					dataField = pieces[0]

					if len(pieces) > 1 {
						for y := 1; y < len(pieces); y++ {
							if pieces[y] == FieldFlagCSV {
								isCSV = true
							} else if pieces[y] == FieldFlagBase64 {
								isBase64 = true
							} else if pieces[y] == FieldFlagBytes {
								isBytes = true
							}
						}
					}
				}
			}
		}

		if isCSV {
			if typed, isTyped := fieldValue.Interface().([]string); isTyped {
				output[dataField] = strings.Join(typed, ",")
			}
		} else if isBytes {
			if typed, isTyped := fieldValue.Interface().([]byte); isTyped {
				output[dataField] = string(typed)
			}
		} else if isBase64 {
			if typed, isTyped := fieldValue.Interface().([]byte); isTyped {
				output[dataField] = base64.StdEncoding.EncodeToString(typed)
			}
			if typed, isTyped := fieldValue.Interface().(string); isTyped {
				output[dataField] = typed
			}
		} else {
			output[dataField] = fmt.Sprintf("%v", FollowValue(fieldValue))
		}
	}

	return output
}
