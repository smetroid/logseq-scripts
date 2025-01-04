package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
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

// Helper function to find the index of the line corresponding to a character offset
func findLineIndex(lines []string, offset int) int {
	currentOffset := 0
	for i, line := range lines {
		currentOffset += len(line) + 1 // Add 1 for the newline character
		if currentOffset > offset {
			return i
		}
	}
	return -1
}

// extractMdBlocks extracts fenced code blocks and their associated tags from markdown text.
func extractMdBlocks(text []byte) []map[string]string {
	// Define the regular expression pattern for fenced code blocks
	codeBlockPattern := regexp.MustCompile("(?s)```(?:\\w+\\s+)?(.*?)```")

	// Find all matches of code blocks in the text
	matches := codeBlockPattern.FindAllStringSubmatchIndex(string(text), -1)

	var results []map[string]string
	lines := strings.Split(string(text), "\n") // Split the text into lines for tag detection

	for _, match := range matches {
		if len(match) >= 4 {
			// Extract the code block
			codeBlock := strings.TrimSpace(string(text[match[2]:match[3]]))

			// Locate the potential tag line
			tags := ""
			endOfBlockIndex := match[1] // End of the matched code block
			tagStartLine := findLineIndex(lines, endOfBlockIndex)

			if tagStartLine >= 0 && tagStartLine+1 < len(lines) {
				possibleTagLine := strings.TrimSpace(lines[tagStartLine+1])
				if strings.HasPrefix(possibleTagLine, "#") { // Ensure it's a valid tag line
					tags = possibleTagLine
				}
			}

			// Add the code block and tag to the results
			results = append(results, map[string]string{
				"code": codeBlock,
				"tags": tags,
			})
		}
	}

	return results
}

// Function to filter commands based on a tag
func filterCommandsByTag(resultSlice []map[string]string, tag string) map[string]string {
	// Create a map to hold the filtered results
	filteredResults := make(map[string]string)

	// Iterate through the slice and check if the "tags" contain the specified tag
	for _, result := range resultSlice {
		// Check if the "tags" contain the provided tag (case-sensitive)
		if strings.Contains(result["tags"], tag) {
			// Add the command to the filtered map
			filteredResults[result["command"]] = result["tags"]
		}
	}

	return filteredResults
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
	results := make(chan map[string]string, len(files))

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
				results <- match
			}
		}(file)
	}

	// Close the results channel once all Go routines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var slice []map[string]string
	for cmd := range results {
		slice = append(slice, cmd)
	}

	// Sort the results by the "command" key (alphabetical order)
	sort.Slice(slice, func(i, j int) bool {
		return slice[i]["command"] < slice[j]["command"]
	})

	// Read all data from the channel into a slice
	stringsMap := []map[string]string{}
	for _, item := range slice {
		itemMap := map[string]string{"command": item["code"], "tags": item["tags"]}
		stringsMap = append(stringsMap, itemMap)
	}

	// Call the function with tag "docker"
	tag := "cmd"
	filteredCommands := filterCommandsByTag(stringsMap, tag)

	// Convert list to JSON
	jsonData, err := json.MarshalIndent(filteredCommands, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	// Print JSON to console
	fmt.Println(string(jsonData))
}
