// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

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
