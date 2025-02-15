package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ConcatenateTxtFiles scans the given directory (and its subdirectories) for .txt files.
// For each file found, it reads the content, extracts the file's parent folder name, and
// then passes both the file content and the folder name to the provided callback function.
// The string returned by the callback is concatenated into the final result.
func ConcatenateTxtFiles(rootPath string, callback func(text, dir string) string) (string, error) {
	var builder strings.Builder

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Return error if there's an issue accessing the file.
			return err
		}

		// Process only files with a ".txt" extension (case-insensitive).
		if !info.IsDir() && strings.EqualFold(filepath.Ext(info.Name()), ".txt") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Get the name of the folder that contains the file.
			// This uses the base name of the parent directory.
			parentDir := filepath.Base(filepath.Dir(path))
			// Process the file content with the callback.
			processed := callback(string(data), parentDir)
			builder.WriteString(processed)
			builder.WriteString("\n")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func readScript() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine user home directory:", err)
	}

	// Construct the path to the IdleYou scripts folder
	scriptPath := filepath.Join(homeDir, "Documents", "IdleYou", "scripts")

	// Make sure all folders exist
	err = os.MkdirAll(scriptPath, 0755)
	if err != nil {
		log.Fatal("Could not create scripts folder:", err)
	}

	// If the scripts folder has script files, use them
	script, err := ConcatenateTxtFiles(scriptPath, func(text string, dir string) string {
		// Replace all lines starting with "=== " with "=== dir/".
		// Replace all lines starting with "! show " with "! show dir/".
		// This keeps mods from overriding other events by prefixing event names
		// with the name of the mod, which is the directory name.
		//
		// Doesn't prefix if it's the root scripts directory.
		var builder strings.Builder
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "=== ") && dir != "scripts" {
				line = fmt.Sprintf("=== %s/%s", dir, line[4:])
			}
			if strings.HasPrefix(line, "! show ") && dir != "scripts" {
				line = fmt.Sprintf("! show %s/%s", dir, line[7:])
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		return builder.String()
	})
	if err != nil {
		log.Printf("Error reading script files: %v\n", err)
	}
	if err == nil && script != "" {
		return script
	}

	// Otherwise, read from the embedded file
	data, err := scriptFile.ReadFile("script.txt")
	if err != nil {
		log.Fatal("Error reading embedded script:", err)
	}

	// and write it to the scripts folder so it can be modified by the user
	err = os.WriteFile(filepath.Join(scriptPath, "script.txt"), data, 0644)
	if err != nil {
		log.Fatal("Error writing script file:", err)
	}

	return string(data)
}
