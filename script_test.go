package main

import (
	"testing"
)

// -----------------------------
// Tests for Utility Functions
// -----------------------------

func TestGetVariableType(t *testing.T) {
	tests := []struct {
		input          interface{}
		expectedInt    bool
		expectedFloat  bool
		expectedString bool
		expectedBool   bool
	}{
		{10, true, false, false, false},
		{3.14, false, true, false, false},
		{"hello", false, false, true, false},
		{true, false, false, false, true}, // Booleans should return all false
	}

	for _, test := range tests {
		isInt, isFloat, isString, isBool := getVariableType(test.input)
		if isInt != test.expectedInt || isFloat != test.expectedFloat || isString != test.expectedString || isBool != test.expectedBool {
			t.Errorf("getVariableType(%v) = (%v, %v, %v, %v), expected (%v, %v, %v, %v)",
				test.input, isInt, isFloat, isString, isBool, test.expectedInt, test.expectedFloat, test.expectedString, test.expectedBool)
		}
	}
}

func TestGetVariableValue(t *testing.T) {
	tests := []struct {
		input          interface{}
		expectedInt    int
		expectedFloat  float64
		expectedString string
		expectedBool   bool
	}{
		{10, 10, 0, "", false},
		{3.14, 0, 3.14, "", false},
		{"hello", 0, 0, "hello", false},
		{true, 0, 0, "", true},
	}

	for _, test := range tests {
		intVal, floatVal, strVal, boolVal := getVariableValue(test.input)
		if intVal != test.expectedInt || floatVal != test.expectedFloat || strVal != test.expectedString || boolVal != test.expectedBool {
			t.Errorf("getVariableValue(%v) = (%d, %f, %s, %t), expected (%d, %f, %s, %t)",
				test.input, intVal, floatVal, strVal, boolVal, test.expectedInt, test.expectedFloat, test.expectedString, test.expectedBool)
		}
	}
}

// -----------------------------
// Tests for Condition Evaluation
// -----------------------------

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		operator string
		varInt   int
		valInt   int
		varFloat float64
		valFloat float64
		expected bool
	}{
		{"<", 5, 10, 0, 0, true},
		{">", 10, 5, 0, 0, true},
		{"==", 10, 10, 0, 0, true},
		{"!=", 10, 5, 0, 0, true},
		{"<=", 5, 5, 0, 0, true},
		{">=", 10, 10, 0, 0, true},
		{"<", 0, 0, 1.5, 2.5, true},
		{">", 0, 0, 2.5, 1.5, true},
	}

	for _, test := range tests {
		result := evaluateCondition(test.operator, test.varInt, test.valInt, test.varFloat, test.valFloat)
		if result != test.expected {
			t.Errorf("evaluateCondition(%q, %d, %d, %f, %f) = %v, expected %v",
				test.operator, test.varInt, test.valInt, test.varFloat, test.valFloat, result, test.expected)
		}
	}
}

func TestEvaluateNumericCondition(t *testing.T) {
	tests := []struct {
		operator  string
		varValue  interface{}
		condValue interface{}
		expected  bool
	}{
		{">", 10, 5, true},
		{"<=", 5, 5, true},
		{"==", 10, 10, true},
		{"!=", 10, 5, true},
		{"<", 3.5, 4.2, true},
		{">=", 7.1, 7.1, true},
	}

	for _, test := range tests {
		result := evaluateNumericCondition(test.operator, test.varValue, test.condValue)
		if result != test.expected {
			t.Errorf("evaluateNumericCondition(%q, %v, %v) = %v, expected %v",
				test.operator, test.varValue, test.condValue, result, test.expected)
		}
	}
}

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
