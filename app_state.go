package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2/data/binding"
)

type AppState struct {
	Work       binding.Int
	Food       binding.Int
	FoodMax    binding.Int
	Mood       binding.Int
	Money      binding.Int
	Job        binding.String
	Working    binding.Bool
	EventName  binding.String
	EventValue binding.Int
	EventMax   binding.Int
	Messages   binding.StringList
	Events     []Event
}

func NewAppState(workValue, foodValue, moodValue, moneyValue int, job string, working bool, eventName string, eventValue int, eventMax int) *AppState {
	var appstate AppState
	appstate = AppState{
		Work:       binding.NewInt(),
		Food:       binding.NewInt(),
		FoodMax:    binding.NewInt(),
		Mood:       binding.NewInt(),
		Money:      binding.NewInt(),
		Job:        binding.NewString(),
		Working:    binding.NewBool(),
		EventName:  binding.NewString(),
		EventValue: binding.NewInt(),
		EventMax:   binding.NewInt(),
		Messages:   binding.NewStringList(),
		Events:     GetEvents(&appstate),
	}
	appstate.Work.Set(workValue)
	appstate.Food.Set(foodValue)
	appstate.FoodMax.Set(foodValue)
	appstate.Mood.Set(moodValue)
	appstate.Money.Set(moneyValue)
	appstate.Job.Set(job)
	appstate.Working.Set(working)
	appstate.EventName.Set(eventName)
	appstate.EventValue.Set(eventValue)
	appstate.EventMax.Set(eventMax)
	return &appstate
}

func NewAppStateWithDefaults() *AppState {
	return NewAppState(
		0,        // workValue
		100,      // foodValue
		50,       // moodValue
		100,      // moneyValue
		"Intern", // job
		true,     // working
		"",       // eventName
		0,        // eventValue
		100,      // eventMax
	)
}

func fromJSON(jsonData string) *AppState {
	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil
	}
	workValue := data["work"]
	foodValue := data["food"]
	moodValue := data["mood"]
	moneyValue := data["money"]
	job := data["job"]
	working := data["working"]
	eventName := data["eventName"]
	eventValue := data["eventValue"]
	eventMax := data["eventMax"]

	return NewAppState(
		workValue.(int),
		foodValue.(int),
		moodValue.(int),
		moneyValue.(int),
		job.(string),
		working.(bool),
		eventName.(string),
		eventValue.(int),
		eventMax.(int),
	)
}

func (state *AppState) gameTick() {
	// Work
	v, err := state.Work.Get()
	if err != nil {
		fmt.Println("Error getting work:", err)
		return
	}
	working, err := state.Working.Get()
	if err != nil {
		fmt.Println("Error getting working:", err)
		return
	}
	if working {
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

	// Handle current event
	eventName, err := state.EventName.Get()
	if err != nil {
		fmt.Println("Error getting event name:", err)
		return
	}
	if eventName != "" {
		eventValue, err := state.EventValue.Get()
		if err != nil {
			fmt.Println("Error getting event value:", err)
			return
		}
		state.EventValue.Set(eventValue + 1)
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
	moodValue, err := state.Mood.Get()
	if err != nil {
		return "", err
	}
	moneyValue, err := state.Money.Get()
	if err != nil {
		return "", err
	}
	job, err := state.Job.Get()
	if err != nil {
		return "", err
	}
	eventName, err := state.EventName.Get()
	if err != nil {
		return "", err
	}
	eventValue, err := state.EventValue.Get()
	if err != nil {
		return "", err
	}
	eventMax, err := state.EventMax.Get()
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(map[string]any{
		"work":       workValue,
		"food":       foodValue,
		"mood":       moodValue,
		"money":      moneyValue,
		"job":        job,
		"eventName":  eventName,
		"eventValue": eventValue,
		"eventMax":   eventMax,
	})
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
