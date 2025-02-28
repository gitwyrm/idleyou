// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

import (
	"testing"

	"fyne.io/fyne/v2/data/binding"
)

// ---------------------------------------
// Tests for fromJSON / toJSON of AppState
// ---------------------------------------

func TestAppStateFromJSONtoJSON(t *testing.T) {
	jsonString, err := NewAppStateWithDefaults().toJSON()
	if err != nil {
		t.Errorf("Error converting AppState to JSON: %s", err)
		return
	}
	appState := fromJSON(jsonString)

	if appState == nil {
		t.Errorf("Expected appState to be non-nil")
		return
	}

	// Check default values
	checkBindingInt(t, appState.Ticks, 0)
	checkBindingInt(t, appState.Work, 0)
	checkBindingInt(t, appState.WorkXP, 0)
	checkBindingInt(t, appState.Food, 200)
	checkBindingInt(t, appState.FoodMax, 200)
	checkBindingInt(t, appState.Energy, 100)
	checkBindingInt(t, appState.EnergyMax, 100)
	checkBindingInt(t, appState.Mood, 50)
	checkBindingInt(t, appState.Money, 100)
	checkBindingInt(t, appState.Charisma, 0)
	checkBindingInt(t, appState.Fitness, 0)
	checkBindingString(t, appState.Job, "")
	checkBindingInt(t, appState.Salary, 0)
	checkBindingBool(t, appState.Working, false)
	checkBindingBool(t, appState.Paused, false)
	checkBindingBool(t, appState.RoutineShower, true)
	checkBindingBool(t, appState.RoutineShave, false)
	checkBindingBool(t, appState.RoutineBrushTeeth, true)
	checkBindingInt(t, appState.RoutineBonus, 0)
	checkBindingString(t, appState.ProgressEventName, "")
	checkBindingInt(t, appState.ProgressEventValue, 0)
	checkBindingInt(t, appState.ProgressEventMax, 100)
	checkBindingString(t, appState.ChoiceEventName, "")
	checkBindingString(t, appState.ChoiceEventText, "")
	checkBindingStringList(t, appState.ChoiceEventChoices, []string{})
	checkBindingStringList(t, appState.Messages, []string{})
	checkBindingUntypedMap(t, appState.Variables, map[string]any{})
}

// Helper functions to check binding values

func checkBindingInt(t *testing.T, b binding.Int, expected int) {
	v, err := b.Get()
	if err != nil {
		t.Errorf("Error getting binding.Int value: %s", err)
		return
	}
	if v != expected {
		t.Errorf("Expected %d, got %d", expected, v)
	}
}

func checkBindingString(t *testing.T, b binding.String, expected string) {
	v, err := b.Get()
	if err != nil {
		t.Errorf("Error getting binding.String value: %s", err)
		return
	}
	if v != expected {
		t.Errorf("Expected %s, got %s", expected, v)
	}
}

func checkBindingBool(t *testing.T, b binding.Bool, expected bool) {
	v, err := b.Get()
	if err != nil {
		t.Errorf("Error getting binding.Bool value: %s", err)
		return
	}
	if v != expected {
		t.Errorf("Expected %t, got %t", expected, v)
	}
}

func checkBindingStringList(t *testing.T, b binding.StringList, expected []string) {
	v, err := b.Get()
	if err != nil {
		t.Errorf("Error getting binding.StringList value: %s", err)
		return
	}
	if len(v) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(v))
		return
	}
	for i := range v {
		if v[i] != expected[i] {
			t.Errorf("Expected %s, got %s at index %d", expected[i], v[i], i)
		}
	}
}

func checkBindingUntypedMap(t *testing.T, b binding.UntypedMap, expected map[string]any) {
	v, err := b.Get()
	if err != nil {
		t.Errorf("Error getting binding.UntypedMap value: %s", err)
		return
	}
	if len(v) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(v))
		return
	}
	for key, value := range expected {
		if v[key] != value {
			t.Errorf("Expected %v, got %v for key %s", value, v[key], key)
		}
	}
}
