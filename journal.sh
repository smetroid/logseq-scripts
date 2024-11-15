#!/bin/bash

LOGSEQ_DIR="$HOME/Documents/logseq"  # Replace with your Logseq graph directory
TODAY=$(date +"%Y_%m_%d")
JOURNAL_FILE="$LOGSEQ_DIR/journals/$TODAY.md"

# Create a new journal file if it doesn't exist
if [ ! -f "$JOURNAL_FILE" ]; then
  echo "---" >> "$JOURNAL_FILE"
  echo "title:: $TODAY" >> "$JOURNAL_FILE"
  echo "---" >> "$JOURNAL_FILE"
  # Append the entry to the journal file
  echo "## Entry from Bash Script" >> "$JOURNAL_FILE"
fi

# Accessing the last command in a different way
CMD=$(HISTFILE=~/.bash_history && history -r && history | tail -n2 | head -n1 | sed 's/^[ ]*[0-9]*[ ]*//')
echo "CMD: $CMD"

# Get the first word of the CMD as a tag
TAG=$(echo "$CMD" | cut -d " " -f1)

# Initialize an array with the default tags
TAGS=("#cmd" "#shell" "#$TAG")

# Prompt for additional tags
read -p "Enter additional tags (space-separated): " -a ADDITIONAL_TAGS

# Add a '#' prefix to each additional tag and append to the TAGS array
for tag in "${ADDITIONAL_TAGS[@]}"; do
  TAGS+=("#$tag")
done

# Join the tags array into a space-separated string
TAGS_STRING="${TAGS[*]}"

# Append the formatted content to the journal file
{
  echo ""
  echo "- $(date):"
  echo '  ```shell'
  echo "  $CMD"
  echo '  ```'
  echo "  $TAGS_STRING"
} >> "$JOURNAL_FILE"
