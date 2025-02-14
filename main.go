package main

import (
	"embed"
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
