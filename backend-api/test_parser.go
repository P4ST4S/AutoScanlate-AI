package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var progressPattern = regexp.MustCompile(`PROGRESS:\s+(\d+)%\s+-\s+(.+)`)

func parseProgressLine(line string) (int, string) {
	if progressPattern.MatchString(line) {
		matches := progressPattern.FindStringSubmatch(line)
		if len(matches) > 2 {
			progress, _ := strconv.Atoi(matches[1])
			message := matches[2]
			return progress, message
		}
	}
	return -1, ""
}

func main() {
	testLines := []string{
		"PROGRESS: 15% - Translated 1/6 pages",
		"PROGRESS: 30% - Translated 2/6 pages",
		"PROGRESS: 50% - Translated 3/6 pages",
		"   Processing: image_01.jpg (1/6)",
		"Some other line",
	}

	for _, line := range testLines {
		progress, message := parseProgressLine(line)
		if progress >= 0 {
			fmt.Printf("✓ Parsed: %d%% - %s\n", progress, message)
		} else {
			fmt.Printf("✗ Not matched: %s\n", line)
		}
	}
}
