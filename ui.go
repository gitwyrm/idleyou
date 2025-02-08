package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func progressBarForBinding(b binding.Int, max binding.Int) *widget.ProgressBar {
	progress := widget.NewProgressBar()
	b.AddListener(binding.NewDataListener(func() {
		v, err := b.Get()
		if err != nil {
			fmt.Println("Error getting work:", err)
			return
		}

		var progressFloat float64
		if max != nil {
			maxValue, err := max.Get()
			if err != nil {
				fmt.Println("Error getting max value:", err)
				return
			}
			progressFloat = float64(v) / float64(maxValue)
		} else {
			progressFloat = float64(v) / 100.0
		}
		progress.SetValue(progressFloat)
	}))
	return progress
}

func setupUI(appstate *AppState) *fyne.Container {

	progressContainer := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Work"), progressBarForBinding(appstate.Work, nil),
		widget.NewLabel("Food"), progressBarForBinding(appstate.Food, appstate.FoodMax),
		widget.NewLabel("Mood"), progressBarForBinding(appstate.Mood, nil),
	)

	// Add a button to save state
	saveButton := widget.NewButton("Save", func() {
		jsonData, err := appstate.toJSON()
		if err != nil {
			fmt.Println("Error saving state:", err)
			return
		}
		fmt.Println("Saved state:", jsonData)
	})

	moneyRow := container.New(layout.NewHBoxLayout(), widget.NewLabel("Money:"), widget.NewLabelWithData(binding.IntToString(appstate.Money)))
	buttonRow := container.New(
		layout.NewHBoxLayout(),
		widget.NewButton("Buy food ($100)", func() {
			money, err := appstate.Money.Get()
			if err != nil {
				fmt.Println("Error getting money:", err)
				return
			}
			food, err := appstate.Food.Get()
			if err != nil {
				fmt.Println("Error getting food:", err)
				return
			}
			if money >= 100 {
				appstate.Money.Set(money - 100)
				appstate.Food.Set(food + 100)
				appstate.FoodMax.Set(food + 100)
				appstate.Messages.Prepend("You bought food!")
			}
		}),
		widget.NewButton("Toggle Work", func() {
			working, err := appstate.Working.Get()
			if err != nil {
				fmt.Println("Error getting working state:", err)
				return
			}
			appstate.Working.Set(!working)
		}))

	messageList := widget.NewListWithData(appstate.Messages,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	// Buttons and progress bars
	column := container.New(layout.NewVBoxLayout(), progressContainer, moneyRow, buttonRow, saveButton)
	// Column with buttons and progress at the top, messages at the center so they fill the space
	return container.NewBorder(column, nil, nil, nil, messageList)
}
