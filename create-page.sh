#!/bin/bash

# Set the title of the new page
title="New Page from Bash"

# Set the content of the new page
content="This is a new page created from a bash script."

# Create the new page file
logseq_dir="$HOME/Documents/logseq"  # Replace with your Logseq graph directory
page_file="$logseq_dir/pages/$title.md"
touch "$page_file"

# Write the content to the page file
echo "# $title" > "$page_file"
echo "$content" >> "$page_file"

echo "New page created: $page_file"
