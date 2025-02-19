// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ConcatenateTxtFiles scans the given directory (and its subdirectories) for .txt files.
// For each file found, it reads the content, extracts the modName, and
// then passes both the file content and the modName to the provided callback function.
// The string returned by the callback is concatenated into the final result.
func ConcatenateTxtFiles(rootPath string, callback func(text, modName string) string) (string, error) {
	var builder strings.Builder

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.EqualFold(filepath.Ext(info.Name()), ".txt") {
			// Only process if the parent dir is "scripts"
			parentDir := filepath.Base(filepath.Dir(path))
			if parentDir != "scripts" {
				return nil // Skip this file
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Mod name is the dir above "scripts"
			modName := filepath.Base(filepath.Dir(filepath.Dir(path)))
			if modName == "mods" { // Root-level check
				modName = ""
			}
			processed := callback(string(data), modName)
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine user home directory:", err)
	}
	modPath := filepath.Join(homeDir, "Documents", "IdleYou", "mods")
	err = os.MkdirAll(modPath, 0755)
	if err != nil {
		log.Fatal("Could not create mods folder:", err)
	}

	script, err := ConcatenateTxtFiles(modPath, func(text, modName string) string {
		var builder strings.Builder
		if modName == "" || modName == "mods" { // Skip root-level prefixing
			return text
		}
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "=== ") {
				line = fmt.Sprintf("=== %s/%s", modName, line[4:])
			}
			if strings.HasPrefix(line, "! show ") {
				line = fmt.Sprintf("! show %s/images/%s", modName, line[7:])
			}
			if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "*") {
				// prefix the eventName after the -> with the modName
				parts := strings.Split(line, "->")
				if len(parts) == 2 {
					line = fmt.Sprintf("%s -> %s/%s", strings.Trim(parts[0], " "), modName, strings.Trim(parts[1], " "))
				}
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		return builder.String()
	})
	if err != nil {
		log.Printf("Error reading mod script files: %v\n", err)
	}
	if script != "" {
		fmt.Println(script)
		return script
	}

	// Fallback to default mod
	data, err := scriptFile.ReadFile("script.txt")
	if err != nil {
		log.Fatal("Error reading embedded script:", err)
	}
	defaultModScriptPath := filepath.Join(modPath, "default", "scripts")
	err = os.MkdirAll(defaultModScriptPath, 0755)
	if err != nil {
		log.Fatal("Could not create default mod scripts folder:", err)
	}
	err = os.WriteFile(filepath.Join(defaultModScriptPath, "script.txt"), data, 0644)
	if err != nil {
		log.Fatal("Error writing default script file:", err)
	}
	return string(data)
}

func getStringAfterSlash(s string) string {
	_, after, found := strings.Cut(s, "/")
	if found {
		return after
	}
	return s
}
