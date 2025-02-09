package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2/data/binding"
)

type AppState struct {
	Work              binding.Int
	WorkXP            binding.Int
	Food              binding.Int
	FoodMax           binding.Int
	Energy            binding.Int
	EnergyMax         binding.Int
	Mood              binding.Int
	Money             binding.Int
	Charisma          binding.Int
	Fitness           binding.Int
	Job               binding.String
	Working           binding.Bool
	RoutineShower     binding.Bool
	RoutineShave      binding.Bool
	RoutineBrushTeeth binding.Bool
	RoutineBonus      binding.Int
	EventName         binding.String
	EventValue        binding.Int
	EventMax          binding.Int
	Messages          binding.StringList
	Events            []Event
}

func NewAppState(workValue, workXP, foodValue, energyValue, energyMaxValue, moodValue, charismaValue, moneyValue, fitnessValue int, job string, working bool, routineShower bool, routineShave bool, routineBrushTeeth bool, routineBonus int, eventName string, eventValue int, eventMax int) *AppState {
	var appstate AppState
	appstate = AppState{
		Work:              binding.NewInt(),
		WorkXP:            binding.NewInt(),
		Food:              binding.NewInt(),
		FoodMax:           binding.NewInt(),
		Energy:            binding.NewInt(),
		EnergyMax:         binding.NewInt(),
		Mood:              binding.NewInt(),
		Money:             binding.NewInt(),
		Charisma:          binding.NewInt(),
		Fitness:           binding.NewInt(),
		Job:               binding.NewString(),
		Working:           binding.NewBool(),
		RoutineShower:     binding.NewBool(),
		RoutineShave:      binding.NewBool(),
		RoutineBrushTeeth: binding.NewBool(),
		RoutineBonus:      binding.NewInt(),
		EventName:         binding.NewString(),
		EventValue:        binding.NewInt(),
		EventMax:          binding.NewInt(),
		Messages:          binding.NewStringList(),
		Events:            GetEvents(&appstate),
	}
	appstate.Work.Set(workValue)
	appstate.WorkXP.Set(workXP)
	appstate.Food.Set(foodValue)
	appstate.FoodMax.Set(foodValue)
	appstate.Energy.Set(energyValue)
	appstate.EnergyMax.Set(energyMaxValue)
	appstate.Mood.Set(moodValue)
	appstate.Money.Set(moneyValue)
	appstate.Charisma.Set(charismaValue)
	appstate.Fitness.Set(fitnessValue)
	appstate.Job.Set(job)
	appstate.Working.Set(working)
	appstate.RoutineShower.Set(routineShower)
	appstate.RoutineShave.Set(routineShave)
	appstate.RoutineBrushTeeth.Set(routineBrushTeeth)
	appstate.RoutineBonus.Set(routineBonus)
	appstate.EventName.Set(eventName)
	appstate.EventValue.Set(eventValue)
	appstate.EventMax.Set(eventMax)
	return &appstate
}

func NewAppStateWithDefaults() *AppState {
	return NewAppState(
		0,        // workValue
		0,        // workXP
		100,      // foodValue
		100,      // energyValue
		100,      // energyMaxValue
		50,       // moodValue
		0,        // charismaValue
		100,      // moneyValue
		0,        // fitnessValue
		"Intern", // job
		true,     // working
		true,     // routineShower
		false,    // routineShave
		true,     // routineBrushTeeth
		0,        // routineBonus
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
	workXP := data["workXP"]
	foodValue := data["food"]
	energyValue := data["energy"]
	energyMaxValue := data["energyMax"]
	moodValue := data["mood"]
	moneyValue := data["money"]
	job := data["job"]
	working := data["working"]
	routineShower := data["routineShower"]
	routineShave := data["routineShave"]
	routineBrushTeeth := data["routineBrushTeeth"]
	routineBonus := data["routineBonus"]
	eventName := data["eventName"]
	eventValue := data["eventValue"]
	eventMax := data["eventMax"]
	charismaValue := data["charisma"]
	fitnessValue := data["fitness"]

	return NewAppState(
		workValue.(int),
		workXP.(int),
		foodValue.(int),
		energyValue.(int),
		energyMaxValue.(int),
		moodValue.(int),
		moneyValue.(int),
		charismaValue.(int),
		fitnessValue.(int),
		job.(string),
		working.(bool),
		routineShower.(bool),
		routineShave.(bool),
		routineBrushTeeth.(bool),
		routineBonus.(int),
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
	workXP, err := state.WorkXP.Get()
	if err != nil {
		fmt.Println("Error getting work XP:", err)
		return
	}
	working, err := state.Working.Get()
	if err != nil {
		fmt.Println("Error getting working:", err)
		return
	}
	eventName, err := state.EventName.Get()
	if err != nil {
		fmt.Println("Error getting event name:", err)
		return
	}

	if working {
		if v < 100 {
			state.Work.Set(v + 1)
			state.WorkXP.Set(workXP + 1)
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

	// Energy
	if working {
		energy, err := state.Energy.Get()
		if err != nil {
			fmt.Println("Error getting energy:", err)
			return
		}
		if energy > 0 {
			mood, err := state.Mood.Get()
			if err != nil {
				fmt.Println("Error getting mood:", err)
				return
			}
			if mood < 50 {
				if energy < 2 {
					state.Energy.Set(0)
				} else {
					state.Energy.Set(energy - 2)
				}
			} else {
				state.Energy.Set(energy - 1)
			}
		} else {
			eventName, err := state.EventName.Get()
			if err != nil {
				fmt.Println("Error getting event name:", err)
				return
			}
			if eventName == "" {
				NewEventHandler(state).Sleep()
			}
		}
	}

	// Food
	v, err = state.Food.Get()
	if err != nil {
		fmt.Println("Error getting food:", err)
		return
	}
	if v > 0 {
		if eventName != "Sleeping" {
			state.Food.Set(v - 1)
		}
	} else {
		if eventName != "Sleeping" {
			fmt.Println("Game Over")
			os.Exit(0)
		}
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
	workXP, err := state.WorkXP.Get()
	if err != nil {
		return "", err
	}
	foodValue, err := state.Food.Get()
	if err != nil {
		return "", err
	}
	energyValue, err := state.Energy.Get()
	if err != nil {
		return "", err
	}
	energyMax, err := state.EnergyMax.Get()
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
	working, err := state.Working.Get()
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
		"workXP":     workXP,
		"food":       foodValue,
		"energy":     energyValue,
		"energyMax":  energyMax,
		"mood":       moodValue,
		"money":      moneyValue,
		"job":        job,
		"working":    working,
		"eventName":  eventName,
		"eventValue": eventValue,
		"eventMax":   eventMax,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (state *AppState) GetAppearance() int {
	fitness, err := state.Fitness.Get()
	if err != nil {
		return 0
	}
	charisma, err := state.Charisma.Get()
	if err != nil {
		return 0
	}
	mood, err := state.Mood.Get()
	if err != nil {
		return 0
	}
	bonus, err := state.RoutineBonus.Get()
	if err != nil {
		return 0
	}
	total := (fitness + charisma + mood) / 3
	total += bonus
	if total > 100 {
		total = 100
	}
	return total
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
