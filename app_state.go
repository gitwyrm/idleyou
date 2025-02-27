// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2/data/binding"
)

type AppState struct {
	// Built-in variables
	Ticks     binding.Int
	Work      binding.Int
	WorkXP    binding.Int
	Food      binding.Int
	FoodMax   binding.Int
	Energy    binding.Int
	EnergyMax binding.Int
	Mood      binding.Int
	Money     binding.Int
	Charisma  binding.Int
	Fitness   binding.Int
	Job       binding.String
	Salary    binding.Int
	Working   binding.Bool
	Paused    binding.Bool
	// Morning routine
	RoutineShower     binding.Bool
	RoutineShave      binding.Bool
	RoutineBrushTeeth binding.Bool
	RoutineBonus      binding.Int
	// Progress events
	ProgressEventName  binding.String
	ProgressEventValue binding.Int
	ProgressEventMax   binding.Int
	// Choice events
	ChoiceEventName    binding.String
	ChoiceEventText    binding.String
	ChoiceEventChoices binding.StringList
	// UI
	Messages binding.StringList
	Events   []Event
	Buttons  binding.UntypedMap
	// Custom variables
	Variables binding.UntypedMap
}

// Adds a persistent button to the UI
func (a *AppState) AddButton(buttonText string, eventName string) {
	// only add button if it doesn't already exist
	if _, err := a.Buttons.GetValue(buttonText); err != nil {
		a.Buttons.SetValue(buttonText, eventName)
	}
}

// Removes a persistent button from the UI
func (a *AppState) RemoveButton(buttonText string) {
	// only remove button if it exists
	if _, err := a.Buttons.GetValue(buttonText); err == nil {
		a.Buttons.Delete(buttonText)
	}
}

// function so set app state variable via string name
// for convenience, variable value will be kept within valid range, for example
// between 0 and 100 for progress bar values
// If the variable is not found, it will be created.
func (a *AppState) Set(variable string, value interface{}) {
	switch strings.ToLower(variable) {
	case "ticks":
		a.Ticks.Set(value.(int))
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
		energyMax, err := a.EnergyMax.Get()
		if err != nil {
			log.Println(err)
			return
		}
		v := value.(int)
		if v < 0 {
			v = 0
		} else if v > energyMax {
			v = energyMax
		}
		a.Energy.Set(v)
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
		a.ProgressEventName.Set(value.(string))
	case "eventvalue":
		a.ProgressEventValue.Set(value.(int))
	case "eventmax":
		a.ProgressEventMax.Set(value.(int))
	case "messages":
		a.Messages.Set(value.([]string))
	case "events":
		a.Events = value.([]Event)
	default:
		a.Variables.SetValue(variable, value)
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
	case "ticks":
		v, err := a.Ticks.Get()
		if err != nil {
			return nil
		}
		return v
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
		v, err := a.ProgressEventName.Get()
		if err != nil {
			return nil
		}
		return v
	case "eventvalue":
		v, err := a.ProgressEventValue.Get()
		if err != nil {
			return nil
		}
		return v
	case "eventmax":
		v, err := a.ProgressEventMax.Get()
		if err != nil {
			return nil
		}
		return v
	default:
		// get the value from Variables
		v, err := a.Variables.GetValue(variable)
		if err != nil {
			return nil
		}
		return v
	}
}

func NewAppState(ticksValue, workValue, workXP, foodValue, foodMaxValue, energyValue, energyMaxValue, moodValue, charismaValue, moneyValue, fitnessValue int, job string, salary int, working bool, paused bool, routineShower bool, routineShave bool, routineBrushTeeth bool, routineBonus int, eventName string, eventValue int, eventMax int, choiceEventName string, choiceEventText string, choiceEventChoices []string, messages []string, variables map[string]any) *AppState {
	appstate := AppState{
		Ticks:              binding.NewInt(),
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
		ProgressEventName:  binding.NewString(),
		ProgressEventValue: binding.NewInt(),
		ProgressEventMax:   binding.NewInt(),
		ChoiceEventName:    binding.NewString(),
		ChoiceEventText:    binding.NewString(),
		ChoiceEventChoices: binding.NewStringList(),
		Messages:           binding.NewStringList(),
		Events:             []Event{},
		Buttons:            binding.NewUntypedMap(),
		Variables:          binding.NewUntypedMap(),
	}
	appstate.Ticks.Set(ticksValue)
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
	appstate.ProgressEventName.Set(eventName)
	appstate.ProgressEventValue.Set(eventValue)
	appstate.ProgressEventMax.Set(eventMax)
	appstate.ChoiceEventName.Set(choiceEventName)
	appstate.ChoiceEventText.Set(choiceEventText)
	appstate.ChoiceEventChoices.Set(choiceEventChoices)
	appstate.Events = GetEvents(&appstate)
	appstate.Messages.Set(messages)
	appstate.Variables.Set(variables)
	return &appstate
}

func NewAppStateWithDefaults() *AppState {
	return NewAppState(
		0,                // ticksValue
		0,                // workValue
		0,                // workXP
		200,              // foodValue
		200,              // foodMaxValue
		100,              // energyValue
		100,              // energyMaxValue
		50,               // moodValue
		0,                // charismaValue
		100,              // moneyValue
		0,                // fitnessValue
		"",               // job
		0,                // salary
		false,            // working
		false,            // paused
		true,             // routineShower
		false,            // routineShave
		true,             // routineBrushTeeth
		0,                // routineBonus
		"",               // eventName
		0,                // eventValue
		100,              // eventMax
		"",               // choiceEventName
		"",               // choiceEventText
		[]string{},       // choiceEventChoices
		[]string{},       // messages
		map[string]any{}, // variables
	)
}

func fromJSON(jsonData string) *AppState {
	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil
	}
	ticksValue := data["ticks"]
	workValue := data["work"]
	workXP := data["workXP"]
	foodValue := data["food"]
	foodMaxValue := data["foodMax"]
	energyValue := data["energy"]
	energyMaxValue := data["energyMax"]
	moodValue := data["mood"]
	moneyValue := data["money"]
	charismaValue := data["charisma"]
	fitnessValue := data["fitness"]
	job := data["job"]
	salary := data["salary"]
	working := data["working"]
	paused := data["paused"]
	routineShower := data["routineShower"]
	routineShave := data["routineShave"]
	routineBrushTeeth := data["routineBrushTeeth"]
	routineBonus := data["routineBonus"]
	progressEventName := data["eventName"]
	progressEventValue := data["eventValue"]
	progressEventMax := data["eventMax"]
	choiceEventName := data["choiceEventName"]
	choiceEventText := data["choiceEventText"]
	choiceEventChoices := data["choiceEventChoices"]
	messages := data["messages"]
	variables := data["variables"]

	return NewAppState(
		ticksValue.(int),
		workValue.(int),
		workXP.(int),
		foodValue.(int),
		foodMaxValue.(int),
		energyValue.(int),
		energyMaxValue.(int),
		moodValue.(int),
		charismaValue.(int),
		moneyValue.(int),
		fitnessValue.(int),
		job.(string),
		salary.(int),
		working.(bool),
		paused.(bool),
		routineShower.(bool),
		routineShave.(bool),
		routineBrushTeeth.(bool),
		routineBonus.(int),
		progressEventName.(string),
		progressEventValue.(int),
		progressEventMax.(int),
		choiceEventName.(string),
		choiceEventText.(string),
		choiceEventChoices.([]string),
		messages.([]string),
		variables.(map[string]any),
	)
}

// Processes a single tick in the game. In other game engines,
// this would be like the update function called in a game loop
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

	// Increment ticks
	ticksValue, err := state.Ticks.Get()
	if err != nil {
		fmt.Println("Error getting ticks:", err)
		return
	}
	ticksValue++
	state.Ticks.Set(ticksValue)

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
	eventName, err := state.ProgressEventName.Get()
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
			eventName, err := state.ProgressEventName.Get()
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

	// Process events in parallel

	// Worker pool setup
	numWorkers := runtime.NumCPU()
	eventQueue := make(chan *Event, len(state.Events))
	var wg sync.WaitGroup

	// Start worker goroutines
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range eventQueue {
				state.handleEvent(event, false)
			}
		}()
	}

	// Enqueue events
	for i := range state.Events {
		eventQueue <- &state.Events[i]
	}
	close(eventQueue)

	wg.Wait()

	// Handle current progress event (if there is one)
	if eventName != "" {
		eventValue, err := state.ProgressEventValue.Get()
		if err != nil {
			fmt.Println("Error getting event value:", err)
			return
		}
		state.ProgressEventValue.Set(eventValue + 1)
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
		for key, value := range event.Choices {
			if value.Conditions != nil {
				// check if all conditions return true
				allTrue := true
				for _, condition := range value.Conditions {
					conditionFn := scriptConditionToFn(state, condition)
					if !conditionFn() {
						allTrue = false
						break
					}
				}
				if !allTrue {
					continue
				}
			}
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
	ticksValue, err := state.Ticks.Get()
	if err != nil {
		return "", err
	}
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
	foodMax, err := state.FoodMax.Get()
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
	charismaValue, err := state.Charisma.Get()
	if err != nil {
		return "", err
	}
	fitnessValue, err := state.Fitness.Get()
	if err != nil {
		return "", err
	}
	job, err := state.Job.Get()
	if err != nil {
		return "", err
	}
	salary, err := state.Salary.Get()
	if err != nil {
		return "", err
	}
	working, err := state.Working.Get()
	if err != nil {
		return "", err
	}
	paused, err := state.Paused.Get()
	if err != nil {
		return "", err
	}
	routineShower, err := state.RoutineShower.Get()
	if err != nil {
		return "", err
	}
	routineShave, err := state.RoutineShave.Get()
	if err != nil {
		return "", err
	}
	routineBrushTeeth, err := state.RoutineBrushTeeth.Get()
	if err != nil {
		return "", err
	}
	routineBonus, err := state.RoutineBonus.Get()
	if err != nil {
		return "", err
	}
	progressEventName, err := state.ProgressEventName.Get()
	if err != nil {
		return "", err
	}
	progressEventValue, err := state.ProgressEventValue.Get()
	if err != nil {
		return "", err
	}
	progressEventMax, err := state.ProgressEventMax.Get()
	if err != nil {
		return "", err
	}
	choiceEventName, err := state.ChoiceEventName.Get()
	if err != nil {
		return "", err
	}
	choiceEventText, err := state.ChoiceEventText.Get()
	if err != nil {
		return "", err
	}
	choiceEventChoices, err := state.ChoiceEventChoices.Get()
	if err != nil {
		return "", err
	}
	messages, err := state.Messages.Get()
	if err != nil {
		return "", err
	}
	variables, err := state.Variables.Get()
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(map[string]any{
		"ticks":              ticksValue,
		"work":               workValue,
		"workXP":             workXP,
		"food":               foodValue,
		"foodMax":            foodMax,
		"energy":             energyValue,
		"energyMax":          energyMax,
		"mood":               moodValue,
		"money":              moneyValue,
		"charisma":           charismaValue,
		"fitness":            fitnessValue,
		"job":                job,
		"salary":             salary,
		"working":            working,
		"paused":             paused,
		"routineShower":      routineShower,
		"routineShave":       routineShave,
		"routineBrushTeeth":  routineBrushTeeth,
		"routineBonus":       routineBonus,
		"progressEventName":  progressEventName,
		"progressEventValue": progressEventValue,
		"progressEventMax":   progressEventMax,
		"choiceEventName":    choiceEventName,
		"choiceEventText":    choiceEventText,
		"choiceEventChoices": choiceEventChoices,
		"messages":           messages,
		"variables":          variables,
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
