// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

import (
	"testing"
)

// -----------------------------
// Tests for Parsing
// -----------------------------

func TestParseCondition(t *testing.T) {
	tests := []struct {
		input    string
		expected ScriptCondition
	}{
		{"mood <= 10", ScriptCondition{"mood", "<=", 10}},
		{"energy > 20", ScriptCondition{"energy", ">", 20}},
		{"status == happy", ScriptCondition{"status", "==", "happy"}},
		{"isRaining == true", ScriptCondition{"isRaining", "==", true}},
		{"routineShave == true", ScriptCondition{"routineShave", "==", true}},
		{"eventName == Morning Routine", ScriptCondition{"eventName", "==", "Morning Routine"}},
	}

	for _, test := range tests {
		result := parseCondition(test.input)
		if result != test.expected {
			t.Errorf("parseCondition(%q) = %+v, expected %+v", test.input, result, test.expected)
		}
	}
}

func TestParseAction(t *testing.T) {
	tests := []struct {
		input    string
		expected ScriptAction
	}{
		{"print You went outside", ScriptAction{"print", "", "You went outside"}},
		{"fitness += 10", ScriptAction{"fitness", "+=", 10}},
		{"energy -= 5", ScriptAction{"energy", "-=", 5}},
	}

	for _, test := range tests {
		result := parseAction(test.input)
		if result != test.expected {
			t.Errorf("parseAction(%q) = %+v, expected %+v", test.input, result, test.expected)
		}
	}
}

func TestParseScript(t *testing.T) {
	script := `=== Go for a walk
? mood <= 10
? energy > 20
! print You went outside for a walk to clear your head.
! fitness += 10
> true`

	scriptEvents := parseScript(script)

	// Check event count
	if len(scriptEvents) != 1 {
		t.Fatalf("Expected 1 script event, got %d", len(scriptEvents))
	}

	scriptEvent := scriptEvents[0]

	// Check event return
	if scriptEvent.Return != true {
		t.Errorf("Expected script event return true, got %v", scriptEvent.Return)
	}

	// Check event name
	if scriptEvent.Name != "Go for a walk" {
		t.Errorf("Expected script event name 'Go for a walk', got '%s'", scriptEvent.Name)
	}

	// Check conditions count
	if len(scriptEvent.ScriptConditions) != 2 {
		t.Fatalf("Expected 2 script conditions, got %d", len(scriptEvent.ScriptConditions))
	}

	// Check actions count
	if len(scriptEvent.ScriptActions) != 2 {
		t.Fatalf("Expected 2 script actions, got %d", len(scriptEvent.ScriptActions))
	}

	// Check first condition
	cond1 := scriptEvent.ScriptConditions[0]
	if cond1.Variable != "mood" || cond1.Operator != "<=" || cond1.Value != 10 {
		t.Errorf("First script condition mismatch: got %+v", cond1)
	}

	// Check second condition
	cond2 := scriptEvent.ScriptConditions[1]
	if cond2.Variable != "energy" || cond2.Operator != ">" || cond2.Value != 20 {
		t.Errorf("Second script condition mismatch: got %+v", cond2)
	}

	// Check first action (print)
	act1 := scriptEvent.ScriptActions[0]
	if act1.Variable != "print" || act1.Value != "You went outside for a walk to clear your head." {
		t.Errorf("First script action mismatch: got %+v", act1)
	}

	// Check second action (fitness += 10)
	act2 := scriptEvent.ScriptActions[1]
	if act2.Variable != "fitness" || act2.Operator != "+=" || act2.Value != 10 {
		t.Errorf("Second script action mismatch: got %+v", act2)
	}
}

// -----------------------------
// Tests for script structs to events
// -----------------------------

func TestScriptEventToEvent(t *testing.T) {
	script := `=== Working a lot
? workxp == 200
! print You worked a lot.
> true`

	scriptEvents := parseScript(script)

	scriptEvent := scriptEvents[0]

	state := NewAppStateWithDefaults()
	state.WorkXP.Set(200)

	event := scriptEventToEvent(state, scriptEvent)

	if !event.Condition() {
		t.Errorf("Expected event condition to be true")
	}

	state.WorkXP.Set(100)

	if event.Condition() {
		t.Errorf("Expected event condition to be false")
	}
}

func TestScriptActionToFn(t *testing.T) {
	scriptAction := ScriptAction{
		Variable: "workxp",
		Operator: "+=",
		Value:    100,
	}

	state := NewAppStateWithDefaults()
	state.WorkXP.Set(100)

	action := scriptActionToFn(state, scriptAction, false)
	action()

	if state.Get("workxp") != 200 {
		t.Errorf("Expected workxp to be 200")
	}
}

func TestScriptConditionToFn(t *testing.T) {
	scriptCondition := ScriptCondition{
		Variable: "workxp",
		Operator: "==",
		Value:    100,
	}

	state := NewAppStateWithDefaults()
	state.WorkXP.Set(100)

	condition := scriptConditionToFn(state, scriptCondition)

	if !condition() {
		t.Errorf("Expected condition to be true")
	}

	scriptCondition.Value = 200

	condition2 := scriptConditionToFn(state, scriptCondition)

	if condition2() {
		t.Errorf("Expected condition2 to be false")
	}

	scriptCondition.Operator = "<"

	condition3 := scriptConditionToFn(state, scriptCondition)

	if !condition3() {
		t.Errorf("Expected condition3 to be true")
	}
}
