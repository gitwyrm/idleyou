package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	// Define constants for game mechanics
	GameSpeed = time.Millisecond * 100
)

//go:embed script.txt
var scriptFile embed.FS

func readScript() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine user home directory:", err)
	}

	// Construct the path to the user's Documents folder
	scriptPath := filepath.Join(homeDir, "Documents", "IdleYou", "script.txt")

	// If a script.txt file exists in the user's documents directory, use that
	if _, err := os.Stat(scriptPath); err == nil {
		data, err := os.ReadFile(scriptPath)
		if err == nil {
			return string(data)
		}
		log.Println("Warning: Error reading script from documents:", err)
	}

	// Otherwise, read from the embedded file
	data, err := scriptFile.ReadFile("script.txt")
	if err != nil {
		log.Fatal("Error reading embedded script:", err)
	}

	return string(data)
}

func main() {
	appstate := NewAppStateWithDefaults()

	a := app.New()
	w := a.NewWindow("IdleYou")

	content := setupUI(appstate)

	appstate.gameTick()

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	go func() {
		for range time.Tick(GameSpeed) {
			appstate.gameTick()
		}
	}()
	w.ShowAndRun()
}
