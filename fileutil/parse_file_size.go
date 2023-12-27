/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package fileutil

import (
	"strconv"
	"strings"
)

// ParseFileSize parses a file size
func ParseFileSize(fileSizeValue string) (int64, error) {
	if len(fileSizeValue) == 0 {
		return 0, nil
	}

	if len(fileSizeValue) < 2 {
		val, err := strconv.Atoi(fileSizeValue)
		if err != nil {
			return 0, err
		}
		return int64(val), nil
	}

	units := strings.ToLower(fileSizeValue[len(fileSizeValue)-2:])
	value, err := strconv.ParseInt(fileSizeValue[:len(fileSizeValue)-2], 10, 64)
	if err != nil {
		return 0, err
	}
	switch units {
	case "tb":
		return value * Terabyte, nil
	case "gb":
		return value * Gigabyte, nil
	case "mb":
		return value * Megabyte, nil
	case "kb":
		return value * Kilobyte, nil
	}
	fullValue, err := strconv.ParseInt(fileSizeValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return fullValue, nil
}
