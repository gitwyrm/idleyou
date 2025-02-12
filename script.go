package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
Simple scripting syntax for events

=== Name of event
? mood <= 10
? energy > 20
! print You went outside for a walk to clear your head.
! fitness += 10
> false

Where lines starting with ? get turned into conditions and
lines starting with ! into actions for the event.
> true/false is the return value of the event, true marks it as done
*/

// -----------------------------
// Script Types
// -----------------------------

type ScriptCondition struct {
	Variable string
	Operator string
	Value    interface{} // string, float64, int, or bool
}

type ScriptAction struct {
	Variable string
	Operator string
	Value    interface{} // string, float64, int, or bool
}

type ScriptEvent struct {
	Name             string
	ScriptConditions []ScriptCondition
	ScriptActions    []ScriptAction
	Choices          map[string]string
	Return           bool
}

// -----------------------------
// Utility Functions
// -----------------------------

// getVariableType returns booleans indicating whether the value is an int, float64, string, or bool
func getVariableType(value interface{}) (isInt, isFloat, isString, isBool bool) {
	switch value.(type) {
	case int:
		return true, false, false, false
	case float64:
		return false, true, false, false
	case string:
		return false, false, true, false
	case bool:
		return false, false, false, true
	default:
		return false, false, false, false
	}
}

// getVariableValue extracts the value as int, float64, string, or bool.
// Only one of the returned values will be non-zero/non-empty based on the type.
func getVariableValue(value interface{}) (int, float64, string, bool) {
	switch v := value.(type) {
	case int:
		return v, 0, "", false
	case float64:
		return 0, v, "", false
	case string:
		return 0, 0, v, false
	case bool:
		return 0, 0, "", v
	default:
		return 0, 0, "", true
	}
}

// evaluateCondition compares two numbers (either both ints or both floats)
// based on the operator.
func evaluateCondition(operator string, varInt, valInt int, varFloat, valFloat float64) bool {
	// If comparing floats, use float comparison
	if varFloat != 0 || valFloat != 0 {
		switch operator {
		case "<":
			return varFloat < valFloat
		case ">":
			return varFloat > valFloat
		case "<=":
			return varFloat <= valFloat
		case ">=":
			return varFloat >= valFloat
		case "==":
			return varFloat == valFloat
		case "!=":
			return varFloat != valFloat
		default:
			panic(fmt.Sprintf("Unknown operator: %s", operator))
		}
	} else {
		// Otherwise, compare ints
		switch operator {
		case "<":
			return varInt < valInt
		case ">":
			return varInt > valInt
		case "<=":
			return varInt <= valInt
		case ">=":
			return varInt >= valInt
		case "==":
			return varInt == valInt
		case "!=":
			return varInt != valInt
		default:
			panic(fmt.Sprintf("Unknown operator: %s", operator))
		}
	}
}

// evaluateNumericCondition determines whether both values are numeric and uses the
// appropriate comparison. If the types do not match, it returns false.
func evaluateNumericCondition(operator string, varVal, condVal interface{}) bool {
	varIsInt, varIsFloat, _, _ := getVariableType(varVal)
	condIsInt, condIsFloat, _, _ := getVariableType(condVal)

	if varIsFloat && condIsFloat {
		_, varFloat, _, _ := getVariableValue(varVal)
		_, condFloat, _, _ := getVariableValue(condVal)
		return evaluateCondition(operator, 0, 0, varFloat, condFloat)
	} else if varIsInt && condIsInt {
		varInt, _, _, _ := getVariableValue(varVal)
		condInt, _, _, _ := getVariableValue(condVal)
		return evaluateCondition(operator, varInt, condInt, 0, 0)
	}
	// Mismatched or non-numeric types result in false
	return false
}

// modifyState applies an action to the state
func modifyState(state *AppState, variable, operator string, value interface{}) {
	currentValue := state.Get(variable)
	switch operator {
	case "=":
		state.Set(variable, value)
	case "+=":
		if intVal, ok := currentValue.(int); ok {
			state.Set(variable, intVal+value.(int))
		} else {
			log.Fatal("Unsupported type for += operation on variable ", variable)
		}
	case "-=":
		if intVal, ok := currentValue.(int); ok {
			state.Set(variable, intVal-value.(int))
		} else {
			log.Fatal("Unsupported type for -= operation on variable ", variable)
		}
	default:
		log.Fatal("Unknown operator: ", operator)
	}
}

// -----------------------------
// Converting ScriptEvent to Event
// -----------------------------

func scriptEventToEvent(state *AppState, scriptEvent ScriptEvent) Event {
	isMultipleChoice := len(scriptEvent.Choices) > 0

	// Build condition functions
	var conditions []func() bool
	for _, condition := range scriptEvent.ScriptConditions {
		conditions = append(conditions, scriptConditionToFn(state, condition))
	}

	// Build action functions
	var actions []func()
	for _, action := range scriptEvent.ScriptActions {
		actions = append(actions, scriptActionToFn(state, action, isMultipleChoice))
	}

	event := NewEvent(
		scriptEvent.Name,
		func() bool {
			for _, condition := range conditions {
				if !condition() {
					return false
				}
			}
			return true
		},
		func() bool {
			for _, action := range actions {
				action()
			}
			return scriptEvent.Return
		},
		scriptEvent.Choices,
	)

	return event
}

// -----------------------------
// Script Action and Condition Functions
// -----------------------------

func scriptActionToFn(state *AppState, action ScriptAction, isMultipleChoice bool) func() {
	// Special case for print commands
	if action.Operator == "" && action.Variable == "print" {
		return func() {
			if isMultipleChoice {
				// TODO: Display message above multiple-choice options
			} else {
				state.Messages.Prepend(action.Value.(string))
			}
		}
	}
	// For other operations, use modifyState
	return func() {
		modifyState(state, action.Variable, action.Operator, action.Value)
	}
}

func scriptConditionToFn(state *AppState, condition ScriptCondition) func() bool {
	// Special case for literal boolean conditions
	if condition.Variable == "boolean" {
		boolVal, ok := condition.Value.(bool)
		if !ok {
			panic("Invalid boolean value in condition")
		}
		return func() bool { return boolVal }
	}

	return func() bool {
		// Retrieve the variable's current value from the state
		variableValue := state.Get(condition.Variable)

		_, _, variableIsString, variableIsBool := getVariableType(variableValue)
		_, _, valueIsString, valueIsBool := getVariableType(condition.Value)

		// If both variable and condition are strings, compare them
		if variableIsString && valueIsString {
			switch condition.Operator {
			case "==":
				return variableValue.(string) == condition.Value.(string)
			case "!=":
				return variableValue.(string) != condition.Value.(string)
			default:
				return false
			}
		}

		// If both variable and condition are booleans, compare them
		if variableIsBool && valueIsBool {
			switch condition.Operator {
			case "==":
				return variableValue.(bool) == condition.Value.(bool)
			case "!=":
				return variableValue.(bool) != condition.Value.(bool)
			default:
				return false
			}
		}

		// Otherwise, treat them as numeric
		return evaluateNumericCondition(condition.Operator, variableValue, condition.Value)
	}
}

// -----------------------------
// Script Parsing Functions
// -----------------------------

func parseScript(script string) []ScriptEvent {
	lines := strings.Split(script, "\n")
	var events []ScriptEvent
	var currentEvent *ScriptEvent

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "==="): // Event name
			if currentEvent != nil {
				events = append(events, *currentEvent)
			}
			currentEvent = &ScriptEvent{
				Name:             strings.TrimSpace(line[3:]),
				ScriptConditions: []ScriptCondition{},
				ScriptActions:    []ScriptAction{},
				Choices:          map[string]string{},
			}

		case strings.HasPrefix(line, "?"): // Condition
			if currentEvent != nil {
				condition := parseCondition(line[1:])
				currentEvent.ScriptConditions = append(currentEvent.ScriptConditions, condition)
			}

		case strings.HasPrefix(line, "!"): // Action
			if currentEvent != nil {
				action := parseAction(line[1:])
				currentEvent.ScriptActions = append(currentEvent.ScriptActions, action)
			}

		case strings.HasPrefix(line, "-"): // Choice
			if currentEvent != nil {
				key, value := parseChoice(line[1:])
				currentEvent.Choices[key] = value
			}

		case strings.HasPrefix(line, ">"):
			if currentEvent != nil {
				retStr := strings.TrimSpace(line[1:])
				switch retStr {
				case "true":
					currentEvent.Return = true
				case "false":
					currentEvent.Return = false
				default:
					log.Fatal("Error parsing ScriptEvent Return:", retStr)
				}
			}
		}
	}

	// Append the final event
	if currentEvent != nil {
		events = append(events, *currentEvent)
	}

	return events
}

// parseChoice parses a choice line in the format: key -> value
func parseChoice(s string) (string, string) {
	parts := strings.Split(s, "->")
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid choice syntax: %s", s))
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value
}

// parseCondition parses a condition line.
// Examples:
//
//	"mood <= 10"  -> variable: mood, operator: <=, value: 10 (int)
//	"status == happy" -> variable: status, operator: ==, value: "happy"
func parseCondition(line string) ScriptCondition {
	parts := strings.Fields(strings.TrimSpace(line))
	// Handle literal booleans
	if len(parts) == 1 {
		var boolean bool
		switch parts[0] {
		case "true":
			boolean = true
		case "false":
			boolean = false
		default:
			panic(fmt.Sprintf("Invalid boolean value: %s", parts[0]))
		}
		return ScriptCondition{
			Variable: "boolean",
			Operator: "",
			Value:    boolean,
		}
	}

	if len(parts) < 3 {
		panic(fmt.Sprintf("Invalid condition syntax: %s", line))
	}

	// Attempt to parse the third part as an int
	if val, err := strconv.Atoi(parts[2]); err == nil {
		return ScriptCondition{
			Variable: parts[0],
			Operator: parts[1],
			Value:    val,
		}
	}

	// Attempt to parse the third part as a float
	if val, err := strconv.ParseFloat(parts[2], 64); err == nil {
		return ScriptCondition{
			Variable: parts[0],
			Operator: parts[1],
			Value:    val,
		}
	}

	// Attempt to parse the third part as a boolean
	if val, err := strconv.ParseBool(parts[2]); err == nil {
		return ScriptCondition{
			Variable: parts[0],
			Operator: parts[1],
			Value:    val,
		}
	}

	// If not an int, bool or float, keep it as a string
	return ScriptCondition{
		Variable: parts[0],
		Operator: parts[1],
		Value:    strings.Join(parts[2:], " "),
	}
}

// parseAction parses an action line.
// Examples:
//
//	"print You went outside" -> variable: print, value: "You went outside"
//	"mood += 10" -> variable: mood, operator: +=, value: 10 (int)
func parseAction(line string) ScriptAction {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) < 2 {
		panic(fmt.Sprintf("Invalid action syntax: %s", line))
	}

	// Special handling for print commands
	if parts[0] == "print" {
		return ScriptAction{
			Variable: "print",
			Operator: "",
			Value:    strings.Join(parts[1:], " "),
		}
	}

	// For other actions, try parsing the value as an int, float64, bool, or string
	if len(parts) >= 3 {
		if val, err := strconv.Atoi(parts[2]); err == nil {
			return ScriptAction{
				Variable: parts[0],
				Operator: parts[1],
				Value:    val,
			}
		}

		if val, err := strconv.ParseBool(parts[2]); err == nil {
			return ScriptAction{
				Variable: parts[0],
				Operator: parts[1],
				Value:    val,
			}
		}

		if val, err := strconv.ParseFloat(parts[2], 64); err == nil {
			return ScriptAction{
				Variable: parts[0],
				Operator: parts[1],
				Value:    val,
			}
		}

		// Otherwise, treat it as a string
		return ScriptAction{
			Variable: parts[0],
			Operator: parts[1],
			Value:    strings.Join(parts[2:], " "),
		}
	}

	panic(fmt.Sprintf("Invalid action syntax: %s", line))
}
