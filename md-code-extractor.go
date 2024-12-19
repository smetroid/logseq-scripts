package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
func extractMdBlocks(text []byte) []string {
	// Define the regular expression pattern for fenced code blocks
	pattern := "```(?:\\w+\\s+)?(.*?)```"

	// Compile the regular expression with the `(?s)` flag to allow . to match newlines
	re := regexp.MustCompile("(?s)" + pattern)

	// Find all matches in the text
	matches := re.FindAllStringSubmatch(string(text), -1)

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
	directory := "~/Documents/logseq/journals/" // Replace with your file path

	// Expand the `~` in the file path
	expandedPath, err := expandPath(directory)
	if err != nil {
		log.Fatalf("Failed to expand path: %v", err)
	}

	// Read all .md files from the directory
	files, err := filepath.Glob(filepath.Join(expandedPath, "*.md"))
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Channel to collect results
	results := make(chan string, len(files))

	// WaitGroup to synchronize Go routines
	var wg sync.WaitGroup

	// Process each file concurrently
	for _, file := range files {
		wg.Add(1) // Increment WaitGroup counter

		go func(file string) {
			defer wg.Done() // Decrement counter when done

			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file, err)
				return
			}

			// Find all code blocks in the file
			matches := extractMdBlocks(content)

			// Collect results
			for _, match := range matches {
				//fmt.Println(match)
				results <- match
			}
		}(file)
	}

	// Close the results channel once all Go routines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Read all data from the channel into a slice
	stringsMap := []map[string]string{}
	for item := range results {
		itemMap := map[string]string{"name": item, "type": "cmd"}
		stringsMap = append(stringsMap, itemMap)
	}

	// Convert list to JSON
	jsonData, err := json.MarshalIndent(stringsMap, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	// Print JSON to console
	fmt.Println(string(jsonData))
}
