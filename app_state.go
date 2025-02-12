package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"fyne.io/fyne/v2/data/binding"
)

type AppState struct {
	Work               binding.Int
	WorkXP             binding.Int
	Food               binding.Int
	FoodMax            binding.Int
	Energy             binding.Int
	EnergyMax          binding.Int
	Mood               binding.Int
	Money              binding.Int
	Charisma           binding.Int
	Fitness            binding.Int
	Job                binding.String
	Salary             binding.Int
	Working            binding.Bool
	Paused             binding.Bool
	RoutineShower      binding.Bool
	RoutineShave       binding.Bool
	RoutineBrushTeeth  binding.Bool
	RoutineBonus       binding.Int
	EventName          binding.String
	EventValue         binding.Int
	EventMax           binding.Int
	ChoiceEventName    binding.String
	ChoiceEventText    binding.String
	ChoiceEventChoices binding.StringList
	Messages           binding.StringList
	Events             []Event
}

// function so set app state variable via string name
// for convenience, variable value will be kept within valid range, for example
// between 0 and 100 for progress bar values
func (a *AppState) Set(variable string, value interface{}) {
	switch strings.ToLower(variable) {
	case "work":
		v := value.(int)
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		a.Work.Set(v)
	case "workxp":
		a.WorkXP.Set(value.(int))
	case "food":
		a.Food.Set(value.(int))
	case "foodmax":
		a.FoodMax.Set(value.(int))
	case "energy":
		a.Energy.Set(value.(int))
	case "energymax":
		a.EnergyMax.Set(value.(int))
	case "mood":
		v := value.(int)
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		a.Mood.Set(v)
	case "money":
		a.Money.Set(value.(int))
	case "charisma":
		v := value.(int)
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		a.Charisma.Set(v)
	case "fitness":
		v := value.(int)
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		a.Fitness.Set(v)
	case "job":
		a.Job.Set(value.(string))
	case "salary":
		a.Salary.Set(value.(int))
	case "working":
		a.Working.Set(value.(bool))
	case "paused":
		a.Paused.Set(value.(bool))
	case "routineshower":
		a.RoutineShower.Set(value.(bool))
	case "routineshave":
		a.RoutineShave.Set(value.(bool))
	case "routinebrushteeth":
		a.RoutineBrushTeeth.Set(value.(bool))
	case "routinebonus":
		a.RoutineBonus.Set(value.(int))
	case "eventname":
		a.EventName.Set(value.(string))
	case "eventvalue":
		a.EventValue.Set(value.(int))
	case "eventmax":
		a.EventMax.Set(value.(int))
	case "messages":
		a.Messages.Set(value.([]string))
	case "events":
		a.Events = value.([]Event)
	}
}

// function to get an Event by name
func (a *AppState) GetEvent(name string) *Event {
	for i, event := range a.Events {
		if event.Name == name {
			return &a.Events[i]
		}
	}
	return nil
}

// function to get app state variable via string name
// to be used with script
func (a *AppState) Get(variable string) interface{} {
	switch strings.ToLower(variable) {
	case "rand":
		// special case that returns random number
		return rand.Float64()
	case "appearance":
		// special case that returns appearance
		return a.GetAppearance()
	case "work":
		v, err := a.Work.Get()
		if err != nil {
			return nil
		}
		return v
	case "workxp":
		v, err := a.WorkXP.Get()
		if err != nil {
			return nil
		}
		return v
	case "food":
		v, err := a.Food.Get()
		if err != nil {
			return nil
		}
		return v
	case "foodmax":
		v, err := a.FoodMax.Get()
		if err != nil {
			return nil
		}
		return v
	case "energy":
		v, err := a.Energy.Get()
		if err != nil {
			return nil
		}
		return v
	case "energymax":
		v, err := a.EnergyMax.Get()
		if err != nil {
			return nil
		}
		return v
	case "mood":
		v, err := a.Mood.Get()
		if err != nil {
			return nil
		}
		return v
	case "money":
		v, err := a.Money.Get()
		if err != nil {
			return nil
		}
		return v
	case "charisma":
		v, err := a.Charisma.Get()
		if err != nil {
			return nil
		}
		return v
	case "fitness":
		v, err := a.Fitness.Get()
		if err != nil {
			return nil
		}
		return v
	case "job":
		v, err := a.Job.Get()
		if err != nil {
			return nil
		}
		return v
	case "salary":
		v, err := a.Salary.Get()
		if err != nil {
			return nil
		}
		return v
	case "working":
		v, err := a.Working.Get()
		if err != nil {
			return nil
		}
		return v
	case "paused":
		v, err := a.Paused.Get()
		if err != nil {
			return nil
		}
		return v
	case "routineshower":
		v, err := a.RoutineShower.Get()
		if err != nil {
			return nil
		}
		return v
	case "routineshave":
		v, err := a.RoutineShave.Get()
		if err != nil {
			return nil
		}
		return v
	case "routinebrushteeth":
		v, err := a.RoutineBrushTeeth.Get()
		if err != nil {
			return nil
		}
		return v
	case "routinebonus":
		v, err := a.RoutineBonus.Get()
		if err != nil {
			return nil
		}
		return v
	case "eventname":
		v, err := a.EventName.Get()
		if err != nil {
			return nil
		}
		return v
	case "eventvalue":
		v, err := a.EventValue.Get()
		if err != nil {
			return nil
		}
		return v
	case "eventmax":
		v, err := a.EventMax.Get()
		if err != nil {
			return nil
		}
		return v
	default:
		return nil
	}
}

func NewAppState(workValue, workXP, foodValue, energyValue, energyMaxValue, moodValue, charismaValue, moneyValue, fitnessValue int, job string, salary int, working bool, paused bool, routineShower bool, routineShave bool, routineBrushTeeth bool, routineBonus int, eventName string, eventValue int, eventMax int, choiceEventName string, choiceEventText string, choiceEventChoices []string) *AppState {
	var appstate AppState
	appstate = AppState{
		Work:               binding.NewInt(),
		WorkXP:             binding.NewInt(),
		Food:               binding.NewInt(),
		FoodMax:            binding.NewInt(),
		Energy:             binding.NewInt(),
		EnergyMax:          binding.NewInt(),
		Mood:               binding.NewInt(),
		Money:              binding.NewInt(),
		Charisma:           binding.NewInt(),
		Fitness:            binding.NewInt(),
		Job:                binding.NewString(),
		Salary:             binding.NewInt(),
		Working:            binding.NewBool(),
		Paused:             binding.NewBool(),
		RoutineShower:      binding.NewBool(),
		RoutineShave:       binding.NewBool(),
		RoutineBrushTeeth:  binding.NewBool(),
		RoutineBonus:       binding.NewInt(),
		EventName:          binding.NewString(),
		EventValue:         binding.NewInt(),
		EventMax:           binding.NewInt(),
		ChoiceEventName:    binding.NewString(),
		ChoiceEventText:    binding.NewString(),
		ChoiceEventChoices: binding.NewStringList(),
		Messages:           binding.NewStringList(),
		Events:             []Event{},
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
	appstate.Salary.Set(salary)
	appstate.Working.Set(working)
	appstate.Paused.Set(paused)
	appstate.RoutineShower.Set(routineShower)
	appstate.RoutineShave.Set(routineShave)
	appstate.RoutineBrushTeeth.Set(routineBrushTeeth)
	appstate.RoutineBonus.Set(routineBonus)
	appstate.EventName.Set(eventName)
	appstate.EventValue.Set(eventValue)
	appstate.EventMax.Set(eventMax)
	appstate.ChoiceEventName.Set(choiceEventName)
	appstate.ChoiceEventText.Set(choiceEventText)
	appstate.ChoiceEventChoices.Set(choiceEventChoices)
	appstate.Events = GetEvents(&appstate)
	return &appstate
}

func NewAppStateWithDefaults() *AppState {
	return NewAppState(
		0,          // workValue
		0,          // workXP
		100,        // foodValue
		100,        // energyValue
		100,        // energyMaxValue
		50,         // moodValue
		0,          // charismaValue
		100,        // moneyValue
		0,          // fitnessValue
		"",         // job
		0,          // salary
		false,      // working
		false,      // paused
		true,       // routineShower
		false,      // routineShave
		true,       // routineBrushTeeth
		0,          // routineBonus
		"",         // eventName
		0,          // eventValue
		100,        // eventMax
		"",         // choiceEventName
		"",         // choiceEventText
		[]string{}, // choiceEventChoices
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
	salary := data["salary"]
	working := data["working"]
	paused := data["paused"]
	routineShower := data["routineShower"]
	routineShave := data["routineShave"]
	routineBrushTeeth := data["routineBrushTeeth"]
	routineBonus := data["routineBonus"]
	eventName := data["eventName"]
	eventValue := data["eventValue"]
	eventMax := data["eventMax"]
	choiceEventName := data["choiceEventName"]
	choiceEventText := data["choiceEventText"]
	choiceEventChoices := data["choiceEventChoices"]
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
		salary.(int),
		working.(bool),
		paused.(bool),
		routineShower.(bool),
		routineShave.(bool),
		routineBrushTeeth.(bool),
		routineBonus.(int),
		eventName.(string),
		eventValue.(int),
		eventMax.(int),
		choiceEventName.(string),
		choiceEventText.(string),
		choiceEventChoices.([]string),
	)
}

func (state *AppState) gameTick() {
	// Pause
	paused, err := state.Paused.Get()
	if err != nil {
		fmt.Println("Error getting paused:", err)
		return
	}
	if paused {
		return
	}

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
			salary, err := state.Salary.Get()
			if err != nil {
				fmt.Println("Error getting salary:", err)
				return
			}
			state.Money.Set(money + salary)
			state.Messages.Prepend(fmt.Sprintf("You were paid $%v for your work!", salary))
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
	for i := range state.Events {
		state.handleEvent(&state.Events[i], false)
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

func (state *AppState) handleEvent(event *Event, ignoreCondition bool) {
	if event.Done {
		return
	}
	if !ignoreCondition && !event.Condition() {
		return
	}
	if len(event.Choices) > 0 {
		keys := make([]string, 0, len(event.Choices)) // Preallocate slice with capacity
		for key := range event.Choices {
			keys = append(keys, key)
		}
		state.ChoiceEventChoices.Set(keys)
		state.ChoiceEventName.Set(event.Name)
	}
	if event.Action() {
		event.Done = true
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
