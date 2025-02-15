// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

type EventHandler struct {
	state *AppState
}

func (e *EventHandler) newHandlerWith(eventName string, doneMessage string, eventMax int, onDone func(), onTick func()) {
	currentEventName, err := e.state.EventName.Get()
	if err != nil {
		fmt.Println("Error getting event name:", err)
		return
	}
	if currentEventName != "" {
		// Already in an event, so just return
		return
	}

	e.state.Working.Set(false)
	e.state.EventName.Set(eventName)
	e.state.EventValue.Set(0)
	e.state.EventMax.Set(eventMax)
	var listener binding.DataListener
	listener = binding.NewDataListener(func() {
		eventValue, err := e.state.EventValue.Get()
		if err != nil {
			fmt.Println("Error getting event value:", err)
			return
		}
		if onTick != nil {
			onTick()
		}
		if eventValue >= eventMax {
			e.state.EventValue.RemoveListener(listener)
			e.state.EventName.Set("")
			e.state.Working.Set(true)
			if doneMessage != "" {
				e.state.Messages.Prepend(doneMessage)
			}
			if onDone != nil {
				onDone()
			}
		}
	})
	e.state.EventValue.AddListener(listener)
}

func (e *EventHandler) MorningRoutine() {
	routineShower, err := e.state.RoutineShower.Get()
	if err != nil {
		fmt.Println("Error getting routine shower:", err)
		return
	}
	routineShave, err := e.state.RoutineShave.Get()
	if err != nil {
		fmt.Println("Error getting routine shave:", err)
		return
	}
	routineBrushTeeth, err := e.state.RoutineBrushTeeth.Get()
	if err != nil {
		fmt.Println("Error getting routine brush teeth:", err)
		return
	}
	ticksNeeded := 0
	bonus := 0
	if routineShower {
		ticksNeeded += 20
		bonus += 10
	}
	if routineShave {
		ticksNeeded += 10
		bonus += 5
	}
	if routineBrushTeeth {
		ticksNeeded += 5
		bonus += 2
	}
	e.newHandlerWith(
		"Morning Routine",
		"You completed your morning routine.",
		ticksNeeded,
		func() {
			e.state.RoutineBonus.Set(bonus)
		},
		nil,
	)
}

func (e *EventHandler) Sleep() {
	e.newHandlerWith(
		"Sleeping",
		"You slept well and feel refreshed.",
		100,
		func() {
			e.MorningRoutine()
		},
		func() {
			energy, err := e.state.Energy.Get()
			if err != nil {
				fmt.Println("Error getting energy:", err)
				return
			}
			energyMax, err := e.state.EnergyMax.Get()
			if err != nil {
				fmt.Println("Error getting energy max:", err)
				return
			}
			if energy < energyMax {
				e.state.Energy.Set(energy + 1)
			}
		},
	)
}

func (e *EventHandler) WatchTV() {
	e.newHandlerWith(
		"Watching TV",
		"You watched TV and feel a little happier, mood increased by 5.",
		100,
		func() {
			mood, err := e.state.Mood.Get()
			if err != nil {
				fmt.Println("Error getting mood:", err)
				return
			}
			if (mood + 5) <= 100 {
				e.state.Mood.Set(mood + 5)
			} else {
				e.state.Mood.Set(100)
			}
		},
		nil,
	)
}

func NewEventHandler(appstate *AppState) *EventHandler {
	return &EventHandler{
		state: appstate,
	}
}
