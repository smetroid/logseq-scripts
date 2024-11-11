package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Set HISTFILE environment variable to point to .bash_history
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	historyFile := homeDir + "/.bash_history"

	// Prepare the command to execute
	cmd := exec.Command("bash", "-c", fmt.Sprintf("HISTFILE=%s history -r && history | tail -n2 | head -n1 | sed 's/^[ ]*[0-9]*[ ]*//'", historyFile))

	// Get output of the command
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	lastCommand := strings.TrimSpace(string(output))
	lastCommand = strings.TrimPrefix(lastCommand, "  ") // Remove leading spaces

	// Define paths and date
	logseqDir := filepath.Join(homeDir, "Documents/logseq") // Adjust as needed
	today := time.Now().Format("2006_01_02")
	journalFile := filepath.Join(logseqDir, "journals", today+".md")

	// Create the journal file if it doesn't exist
	if _, err := os.Stat(journalFile); os.IsNotExist(err) {
		file, err := os.Create(journalFile)
		if err != nil {
			log.Fatalf("Failed to create journal file: %v", err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		fmt.Fprintln(writer, "---")
		fmt.Fprintf(writer, "title:: %s\n", today)
		fmt.Fprintln(writer, "---")
		fmt.Fprintln(writer, "## Entry from Go Script")
		writer.Flush()
	}

	// Get the last command from bash history
	fmt.Printf("CMD: %s\n", lastCommand)

	// Extract the first word as a tag
	tag := strings.Split(lastCommand, " ")[0]

	// Initialize tags with default values
	tags := []string{"#cmd", "#shell", fmt.Sprintf("#%s", tag)}

	// Prompt for additional tags
	fmt.Print("Enter additional tags (space-separated): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		additionalTags := strings.Split(input, " ")
		for _, t := range additionalTags {
			tags = append(tags, fmt.Sprintf("#%s", t))
		}
	}

	// Join tags into a space-separated string
	tagsString := strings.Join(tags, " ")

	// Append the formatted content to the journal file
	file, err := os.OpenFile(journalFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open journal file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "")
	fmt.Fprintf(writer, "- %s:\n", time.Now().Format(time.RFC1123))
	fmt.Fprintln(writer, "  ```shell")
	fmt.Fprintf(writer, "  %s\n", lastCommand)
	fmt.Fprintln(writer, "  ```")
	fmt.Fprintf(writer, "  %s\n", tagsString)
	writer.Flush()

}
