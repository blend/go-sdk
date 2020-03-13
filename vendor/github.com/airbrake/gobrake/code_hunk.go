package gobrake

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/airbrake/gobrake/internal/lrucache"
)

var cache = lrucache.New(1000)

func getCode(file string, line int) (map[int]string, error) {
	cacheKey := file + strconv.Itoa(line)

	v, ok := cache.Get(cacheKey)
	if ok {
		switch v := v.(type) {
		case error:
			return nil, v
		case map[int]string:
			return v, nil
		default:
			return nil, fmt.Errorf("unsupported type=%T", v)
		}
	}

	lines, err := _getCode(file, line)
	if err != nil {
		cache.Set(cacheKey, err)
		return nil, err
	}

	return lines, nil
}

func _getCode(file string, line int) (map[int]string, error) {
	const nlines = 2
	const maxLineLen = 512

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	start := line - nlines
	end := line + nlines
	scanner := bufio.NewScanner(fd)

	var i int
	lines := make(map[int]string, nlines)
	for scanner.Scan() {
		i++
		if i < start {
			continue
		}
		if i > end {
			break
		}
		line := scanner.Text()
		if len(line) > maxLineLen {
			line = line[:maxLineLen]
		}
		lines[i] = line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
