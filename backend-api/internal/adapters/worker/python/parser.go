package python

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regex patterns for parsing Python worker output
	processingPattern = regexp.MustCompile(`Processing:\s+(.+)`)
	pageCountPattern  = regexp.MustCompile(`(\d+)\s+images?`)
	completedPattern  = regexp.MustCompile(`Created:\s+(.+)`)
	errorPattern      = regexp.MustCompile(`❌\s+(.+)`)
)

// parseProgressLine parses a line from the worker's stdout and extracts progress information
// Returns (progress percentage, message)
func parseProgressLine(line string) (int, string) {
	// Check for completion message: "✅ Created: chapter_translated.zip"
	if strings.Contains(line, "✅") && completedPattern.MatchString(line) {
		matches := completedPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			return 100, "Translation completed: " + matches[1]
		}
		return 100, "Translation completed"
	}

	// Check for processing message: "   Processing: image_01.jpg"
	if processingPattern.MatchString(line) {
		matches := processingPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			filename := strings.TrimSpace(matches[1])
			return -1, "Processing: " + filename
		}
	}

	// Check for processed count: " -> Processed 15 images."
	if strings.Contains(line, "Processed") && pageCountPattern.MatchString(line) {
		matches := pageCountPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			count, _ := strconv.Atoi(matches[1])
			return -1, "Processed " + strconv.Itoa(count) + " images"
		}
	}

	// Check for initialization: "⚙️ Device: CUDA"
	if strings.Contains(line, "⚙️ Device:") {
		return 0, "Initializing AI worker"
	}

	// Check for pipeline ready: "✅ Pipeline Ready"
	if strings.Contains(line, "✅ Pipeline Ready") {
		return 5, "AI pipeline ready"
	}

	// Check for errors
	if errorPattern.MatchString(line) {
		matches := errorPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			return -1, "Error: " + matches[1]
		}
	}

	// No parseable information
	return -1, ""
}

// estimateProgress calculates progress percentage based on context
// This is a helper function that can be used when we have more information
// about total pages to process
func estimateProgress(currentPage, totalPages int) int {
	if totalPages == 0 {
		return 0
	}
	progress := (currentPage * 100) / totalPages
	if progress > 100 {
		return 100
	}
	return progress
}
