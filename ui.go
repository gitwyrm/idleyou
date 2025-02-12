package main

import (
	"fmt"
	"log"
	"strings"

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
	// Convert appearance fn to binding for progress bar
	appearanceBinding := binding.NewInt()
	appearanceListener := binding.NewDataListener(func() {
		appearance := appstate.GetAppearance()
		appearanceBinding.Set(appearance)
	})
	appstate.Charisma.AddListener(appearanceListener)
	appstate.Mood.AddListener(appearanceListener)
	appstate.Fitness.AddListener(appearanceListener)
	appstate.RoutineBonus.AddListener(appearanceListener)

	progressContainer := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Food"), progressBarForBinding(appstate.Food, appstate.FoodMax),
		widget.NewLabel("Work"), progressBarForBinding(appstate.Work, nil),
		widget.NewLabel("Energy"), progressBarForBinding(appstate.Energy, appstate.EnergyMax),
		widget.NewLabel("Mood"), progressBarForBinding(appstate.Mood, nil),
		widget.NewLabel("Appearance"), progressBarForBinding(appearanceBinding, nil),
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

	// Choice event buttons
	choiceButtons := container.New(
		layout.NewVBoxLayout(),
	)

	// Choice event container
	choiceContainer := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabelWithData(appstate.ChoiceEventName),
		widget.NewLabelWithData(appstate.ChoiceEventText),
		choiceButtons,
	)

	appstate.ChoiceEventName.AddListener(binding.NewDataListener(func() {
		choiceEventName, err := appstate.ChoiceEventName.Get()
		if err != nil {
			fmt.Println("Error getting choice event name:", err)
			return
		}
		if choiceEventName == "" {
			choiceContainer.Hide()
		} else {
			choiceContainer.Show()
		}
	}))

	appstate.ChoiceEventChoices.AddListener(binding.NewDataListener(func() {
		choiceEventChoices, err := appstate.ChoiceEventChoices.Get()
		if err != nil {
			fmt.Println("Error getting choice event choices:", err)
			return
		}

		// create a button for each choice
		choiceButtons.RemoveAll()
		for _, choice := range choiceEventChoices {
			button := widget.NewButton(choice, func() {
				currentEventName, err := appstate.ChoiceEventName.Get()
				if err != nil {
					fmt.Println("Error getting current event name:", err)
					return
				}
				currentEvent := appstate.GetEvent(currentEventName)
				appstate.ChoiceEventChoices.Set([]string{})
				appstate.ChoiceEventName.Set("")
				event := appstate.GetEvent(currentEvent.Choices[choice])
				if event == nil {
					log.Fatalf("Event not found: '%s'", currentEvent.Choices[choice])
					return
				}
				appstate.handleEvent(event, true)
			})
			choiceButtons.Add(button)
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
		widget.NewButton("Buy food (Max)", func() {
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
			ableToPurchase := money / 100
			if ableToPurchase > 0 {
				appstate.Money.Set(money - ableToPurchase*100)
				appstate.Food.Set(food + ableToPurchase*100)
				appstate.FoodMax.Set(food + ableToPurchase*100)
				appstate.Messages.Prepend(fmt.Sprintf("You bought %v food!", ableToPurchase))
			}
		}),
		widget.NewButton("Watch TV", func() {
			NewEventHandler(appstate).WatchTV()
		}))

	messageList := widget.NewListWithData(appstate.Messages,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			text, err := i.(binding.String).Get()
			if err != nil {
				fmt.Println("Error getting message:", err)
				return
			}

			if strings.HasPrefix(text, "Event: ") {
				text = strings.TrimPrefix(text, "Event: ")
				o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
			} else {
				o.(*widget.Label).TextStyle = fyne.TextStyle{}
			}

			o.(*widget.Label).SetText(text)
		})

	playerInfo := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Job:"), widget.NewLabelWithData(appstate.Job),
		widget.NewLabel("Job experience:"), widget.NewLabelWithData(binding.IntToString(appstate.WorkXP)),
		widget.NewLabel("Money:"), widget.NewLabelWithData(binding.IntToString(appstate.Money)),
		widget.NewLabel("Morning routine"), container.NewHBox(
			widget.NewCheckWithData("Shower", appstate.RoutineShower),
			widget.NewCheckWithData("Shave", appstate.RoutineShave),
			widget.NewCheckWithData("Brush teeth", appstate.RoutineBrushTeeth),
		),
	)

	sidePanel := container.New(layout.NewVBoxLayout(), playerInfo, eventContainer)

	// Buttons and progress bars
	column := container.New(layout.NewVBoxLayout(), progressContainer, buttonRow, saveButton)

	content := container.New(layout.NewGridLayout(3), column, sidePanel, choiceContainer)

	// Column with buttons and progress at the top, messages at the center so they fill the space
	return container.NewBorder(content, nil, nil, nil, messageList)
}
