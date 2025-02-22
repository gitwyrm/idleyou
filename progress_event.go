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

type ProgressEvent struct {
	state *AppState
}

func (e *ProgressEvent) newEventWith(eventName string, doneMessage string, eventMax int, onDone func(), onTick func()) {
	currentEventName, err := e.state.ProgressEventName.Get()
	if err != nil {
		fmt.Println("Error getting event name:", err)
		return
	}
	if currentEventName != "" {
		// Already in an event, so just return
		return
	}

	e.state.Working.Set(false)
	e.state.ProgressEventName.Set(eventName)
	e.state.ProgressEventValue.Set(0)
	e.state.ProgressEventMax.Set(eventMax)
	var listener binding.DataListener
	listener = binding.NewDataListener(func() {
		eventValue, err := e.state.ProgressEventValue.Get()
		if err != nil {
			fmt.Println("Error getting event value:", err)
			return
		}
		if onTick != nil {
			onTick()
		}
		if eventValue >= eventMax {
			e.state.ProgressEventValue.RemoveListener(listener)
			e.state.ProgressEventName.Set("")
			e.state.Working.Set(true)
			if doneMessage != "" {
				e.state.Messages.Prepend(doneMessage)
			}
			if onDone != nil {
				onDone()
			}
		}
	})
	e.state.ProgressEventValue.AddListener(listener)
}

func (e *ProgressEvent) MorningRoutine() {
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
	e.newEventWith(
		"Morning Routine",
		"You completed your morning routine.",
		ticksNeeded,
		func() {
			e.state.RoutineBonus.Set(bonus)
		},
		nil,
	)
}

func (e *ProgressEvent) Sleep() {
	e.newEventWith(
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

func NewEventHandler(appstate *AppState) *ProgressEvent {
	return &ProgressEvent{
		state: appstate,
	}
}
