#!/usr/bin/env bash

# Script to concatenate all project files into a single file

echo -e "# Format Explanation:\n# Each file is prefixed with {{filename.ext}}, followed by its contents.\n# Example:\n# {{example.go}}\n# content of example.go\n" > concatenated.txt
for f in *.mod *.go *.md *.txt; do [[ "$f" != "concatenated.txt" && "$f" != "todo.md" ]] && echo -e "\n{{${f}}}\n" && cat "$f"; done >> concatenated.txt
