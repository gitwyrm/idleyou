package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2/data/binding"
)

type AppState struct {
	Work     binding.Int
	Food     binding.Int
	Mood     binding.Int
	Money    binding.Int
	Job      binding.String
	Messages binding.StringList
	Events   []Event
}

func NewAppState(workValue, foodValue, moodValue, moneyValue int, job string) *AppState {
	var appstate AppState
	appstate = AppState{
		Work:     binding.NewInt(),
		Food:     binding.NewInt(),
		Mood:     binding.NewInt(),
		Money:    binding.NewInt(),
		Job:      binding.NewString(),
		Messages: binding.NewStringList(),
		Events:   GetEvents(&appstate),
	}
	appstate.Work.Set(workValue)
	appstate.Food.Set(foodValue)
	appstate.Mood.Set(moodValue)
	appstate.Money.Set(moneyValue)
	appstate.Job.Set(job)
	return &appstate
}

func NewAppStateWithDefaults() *AppState {
	return NewAppState(0, 100, 50, 100, "Intern")
}

func fromJSON(jsonData string) *AppState {
	var data map[string]int
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil
	}
	workValue := data["work"]
	foodValue := data["food"]
	moodValue := data["mood"]
	moneyValue := data["money"]
	var appstate AppState
	appstate = AppState{
		Work:     binding.NewInt(),
		Food:     binding.NewInt(),
		Mood:     binding.NewInt(),
		Money:    binding.NewInt(),
		Messages: binding.NewStringList(),
		Events:   GetEvents(&appstate),
	}
	appstate.Work.Set(workValue)
	appstate.Food.Set(foodValue)
	appstate.Mood.Set(moodValue)
	appstate.Money.Set(moneyValue)
	return &appstate
}

func (state *AppState) gameTick() {
	// Work
	v, err := state.Work.Get()
	if err != nil {
		fmt.Println("Error getting work:", err)
		return
	}
	if v < 100 {
		state.Work.Set(v + 1)
	} else {
		state.Work.Set(0)
		money, err := state.Money.Get()
		if err != nil {
			fmt.Println("Error getting money:", err)
			return
		}
		state.Money.Set(money + state.GetSalary())
		state.Messages.Prepend(fmt.Sprintf("You were paid $%v for your work!", state.GetSalary()))
	}

	// Food
	v, err = state.Food.Get()
	if err != nil {
		fmt.Println("Error getting food:", err)
		return
	}
	if v > 0 {
		state.Food.Set(v - 1)
	} else {
		fmt.Println("Game Over")
		os.Exit(0)
	}

	// Handle events
	for i, event := range state.Events {
		if event.Done {
			continue
		}
		if !event.Condition() {
			continue
		}
		if event.Action() {
			state.Events[i].Done = true
		}
	}
}

func (state *AppState) toJSON() (string, error) {
	workValue, err := state.Work.Get()
	if err != nil {
		return "", err
	}
	foodValue, err := state.Food.Get()
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(map[string]int{"work": workValue, "food": foodValue})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (state *AppState) GetSalary() int {
	job, err := state.Job.Get()
	if err != nil {
		fmt.Println("Error getting job:", err)
		return 0
	}
	switch job {
	case "Intern":
		return 100
	case "Sales clerk":
		return 200
	case "Manager":
		return 600
	default:
		return 0
	}
}
