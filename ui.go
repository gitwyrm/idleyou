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

	eventContainer := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabelWithData(appstate.EventName),
		progressBarForBinding(appstate.EventValue, appstate.EventMax),
	)

	appstate.EventName.AddListener(binding.NewDataListener(func() {
		eventName, err := appstate.EventName.Get()
		if err != nil {
			fmt.Println("Error getting event name:", err)
			return
		}
		if eventName == "" {
			eventContainer.Hide()
		} else {
			eventContainer.Show()
		}
	}))

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
		widget.NewButton("Watch TV", func() {
			eventName, err := appstate.EventName.Get()
			if err != nil {
				fmt.Println("Error getting event name:", err)
				return
			}
			if eventName != "" {
				// Already in an event, so just return
				return
			}

			appstate.Working.Set(false)
			appstate.EventName.Set("Watching TV")
			appstate.EventValue.Set(0)
			appstate.EventMax.Set(100)
			var listener binding.DataListener
			listener = binding.NewDataListener(func() {
				eventValue, err := appstate.EventValue.Get()
				if err != nil {
					fmt.Println("Error getting event value:", err)
					return
				}
				if eventValue >= 100 {
					mood, err := appstate.Mood.Get()
					if err != nil {
						fmt.Println("Error getting mood:", err)
						return
					}
					appstate.EventName.Set("")
					if (mood + 5) <= 100 {
						appstate.Mood.Set(mood + 5)
					} else {
						appstate.Mood.Set(100)
					}
					appstate.EventValue.RemoveListener(listener)
					appstate.Working.Set(true)
					appstate.Messages.Prepend("You watched TV and feel a little happier, mood increased by 5.")
				}
			})
			appstate.EventValue.AddListener(listener)
		}))

	messageList := widget.NewListWithData(appstate.Messages,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	// Buttons and progress bars
	column := container.New(layout.NewVBoxLayout(), progressContainer, moneyRow, eventContainer, buttonRow, saveButton)
	// Column with buttons and progress at the top, messages at the center so they fill the space
	return container.NewBorder(column, nil, nil, nil, messageList)
}
