#!/bin/bash

set -eu

if [ $# -eq 0 ]; then
    jq -n '{
        title: "Logseq",
        description: "Search logseq code blocks",
        commands: [
            {
                name: "logseq-cmds",
                title: "Logseq cmds blocks",
                mode: "filter"
            },
            {
                name: "logseq-script",
                title: "Logseq script blocks",
                mode: "filter"
            },
            {
                name: "logseq-notes",
                title: "Logseq notes",
                mode: "filter"
            },
            {
                name: "run-command",
                title: "execute command",
                mode: "tty",
                exit: "true"
            },
            {
                name: "view-command",
                title: "view command",
                mode: "detail",
                exit: "false"
            }
        ],
    }'
    exit 0
fi

COMMAND=$(echo "$1" | jq -r '.command')

if [ "$COMMAND" = "logseq-cmds" ]; then
  ~/projects/logseq-scripts/md-code-extractor | jq '{
        "items": map({
            "title": .name,
            "subtitle": .type,
            "actions": [{
                "type": "run",
                "title": "Run cmd ",
                "command": "run-command",
                "params": {
                    "exec": .name,
                }
                },{
                  "type": "run",
                  "title": "view cmd ",
                  "command": "view-command",
                  "params": {
                      "exec": .name,
                  },
            }]
        }),
        "actions": [{
          "title": "Refresh items",
          "type": "reload",
          "exit": "true",
      }]
  }'
  exit 0
fi


if [ "$COMMAND" = "run-command" ]; then
  CMD=$(echo "$1"| jq -r '.params.exec')
  konsole -e bash -c "$CMD; exec bash"
elif [ "$COMMAND" = "view-command" ]; then
    cmd=$(echo "$1"| jq -r '.params.exec')
    jq -n --arg cmd "$cmd" '{
        "text": $cmd,
        "actions": [{
            title: "Copy to clipboard",
            type: "copy",
            text: $cmd,
            exit: false,
        }],
    }'
fi
