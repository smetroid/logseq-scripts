package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		homeDir := usr.HomeDir
		// Replace `~` with the home directory
		path = filepath.Join(homeDir, path[1:])
	}
	return path, nil
}

// extractMdBlocks extracts fenced code blocks from markdown text.
func extractMdBlocks(text string) []string {
	// Define the regular expression pattern for fenced code blocks
	pattern := "```(?:\\w+\\s+)?(.*?)```"

	// Compile the regular expression with the `(?s)` flag to allow . to match newlines
	re := regexp.MustCompile("(?s)" + pattern)

	// Find all matches in the text
	matches := re.FindAllStringSubmatch(text, -1)

	// Extract and trim the matched blocks
	var blocks []string
	for _, match := range matches {
		if len(match) > 1 {
			blocks = append(blocks, strings.TrimSpace(match[1]))
		}
	}

	return blocks
}

func main() {
	// Example file path using `~`
	filePath := "~/Documents/logseq/journals/2024_11_10.md" // Replace with your file path

	// Expand the `~` in the file path
	expandedPath, err := expandPath(filePath)
	if err != nil {
		log.Fatalf("Failed to expand path: %v", err)
	}

	// Read the markdown content from the expanded file path
	content, err := os.ReadFile(expandedPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert the content to a string
	text := string(content)

	// Extract code blocks
	blocks := extractMdBlocks(text)
	for _, block := range blocks {
		fmt.Printf("%s\n", block)
	}
}
