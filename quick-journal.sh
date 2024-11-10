#!/bin/bash

# Get today's date in Logseq format
date=$(date +"%Y-%m-%d")

# Get the journal entry from the command line
read -p "Enter your journal entry: " entry

# Append the entry to the appropriate journal file
logseq_dir="$HOME/Documents/logseq/journals"  # Replace with your actual Logseq journals directory
echo "- $entry" >> "$logseq_dir/$date.md"
