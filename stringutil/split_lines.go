package stringutil

// SplitLines splits a corpus into individual lines by end of line character(s).
// Possible end of line characters include `\n`, `\r` and `\r\n`.
func SplitLines(contents string) []string {
	contentRunes := []rune(contents)

	var output []string

	const newline = '\n'

	var line []rune
	var c rune
	for index := 0; index < len(contentRunes); index++ {
		c = contentRunes[index]
		if c == newline {
			line = append(line, c)
			output = append(output, string(line))
			line = nil
			continue
		}
		line = append(line, c)
		continue
	}

	if len(line) > 0 {
		output = append(output, string(line))
	}

	return output
}
