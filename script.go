// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

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
	ScriptButtons    []ScriptButton
	Choices          map[string]Choice
	ProgressMax      int
	Return           bool
}

// ScriptButton represents a button addition/removal in a ScriptEvent.
// If an EventName is provided, it adds a button.
// If an EventName is not provided, it removes the button with the given ButtonText.
type ScriptButton struct {
	ButtonText string
	EventName  string
}

type Choice struct {
	ButtonText string
	EventName  string
	Conditions []ScriptCondition
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

	// append button addition and removals to actions
	for _, button := range scriptEvent.ScriptButtons {
		if button.EventName == "" {
			actions = append(actions, func() {
				state.RemoveButton(button.ButtonText)
			})
		} else {
			actions = append(actions, func() {
				state.AddButton(button.ButtonText, button.EventName)
			})
		}
	}

	event := NewEvent(
		scriptEvent.Name,
		func() bool {
			// if no conditions are provided, return false so event
			// is never triggered automatically
			if len(conditions) == 0 {
				return false
			}

			for _, condition := range conditions {
				if !condition() {
					return false
				}
			}
			return true
		},
		func() bool {
			// if it is a progress event
			if scriptEvent.ProgressMax > 0 {
				// if no other event is running
				if state.Get("eventName") == "" {
					NewEventHandler(state).newEventWith(
						scriptEvent.Name,
						"",
						scriptEvent.ProgressMax,
						func() {
							for _, action := range actions {
								action()
							}
						},
						nil,
					)
					return scriptEvent.Return
				} else {
					// an event is already running, so return false so that the event
					// can be triggered again even if it is a one off event
					return false
				}
			}

			// if it's not a progress event
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
				state.ChoiceEventText.Set(action.Value.(string))
			} else {
				state.Messages.Prepend(action.Value.(string))
			}
		}
	}

	// Special case for show commands
	if action.Operator == "" && action.Variable == "show" {
		return func() {
			if isMultipleChoice {
				state.ChoiceEventText.Set(action.Value.(string))
			} else {
				state.Messages.Prepend("Image: " + action.Value.(string))
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

		// Skip comments
		if strings.HasPrefix(line, "#") {
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
				Choices:          map[string]Choice{},
			}

		case strings.HasPrefix(line, "?"): // Condition
			if currentEvent != nil {
				condition := parseCondition(line[1:])
				currentEvent.ScriptConditions = append(currentEvent.ScriptConditions, condition)
			}

		case strings.HasPrefix(line, "%"): // ProgressMax
			progressMax := parseProgressMax(line[1:])
			currentEvent.ProgressMax = progressMax

		case strings.HasPrefix(line, "!"): // Action
			if currentEvent != nil {
				action := parseAction(line[1:])
				currentEvent.ScriptActions = append(currentEvent.ScriptActions, action)
			}

		case strings.HasPrefix(line, "*"): // Choice
			if currentEvent != nil {
				key, value, conditions, err := parseChoice(line[1:])
				if err != nil {
					log.Fatal(err)
				}
				currentEvent.Choices[key] = Choice{
					ButtonText: key,
					EventName:  value,
					Conditions: conditions,
				}
			}

		case strings.HasPrefix(line, "+"): // Button Addition
			if currentEvent != nil {
				button := parseButton(line[1:])
				currentEvent.ScriptButtons = append(currentEvent.ScriptButtons, button)
			}

		case strings.HasPrefix(line, "-"): // Button Removal
			if currentEvent != nil {
				button := parseButton(line[1:])
				currentEvent.ScriptButtons = append(currentEvent.ScriptButtons, button)
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

// Parses a line representing a button addition or removal into a ScriptButton
func parseButton(s string) ScriptButton {
	// if the line is in the format: button name -> event name
	// it is a button addition
	if strings.Contains(s, "->") {
		parts := strings.Split(s, "->")
		if len(parts) != 2 {
			log.Fatal("Error parsing ScriptButton:", s)
		}
		return ScriptButton{
			strings.TrimSpace(parts[0]),
			strings.TrimSpace(parts[1]),
		}
	}

	// if the line is in the format: button name
	// it is a button removal
	return ScriptButton{
		strings.TrimSpace(s),
		"",
	}
}

func parseProgressMax(s string) int {
	max, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		log.Printf("Error parsing ScriptEvent ProgressMax: %s\n", s)
		return 0
	}
	return max
}

// parseChoice parses a choice line in the format:
// "variableName < 10, otherVariable == 5: buttonString -> eventName"
// Conditions are optional.
func parseChoice(s string) (string, string, []ScriptCondition, error) {
	parts := strings.SplitN(s, ":", 2) // Split at most once
	var conditions []ScriptCondition

	if len(parts) == 2 {
		conditionStrings := strings.Split(parts[0], ",")
		for _, conditionString := range conditionStrings {
			condition := parseCondition(strings.TrimSpace(conditionString))
			conditions = append(conditions, condition)
		}
		s = parts[1] // Keep only the right-hand side for further parsing
	}

	choiceParts := strings.SplitN(s, "->", 2) // Split at most once
	if len(choiceParts) != 2 {
		return "", "", nil, fmt.Errorf("invalid choice syntax: %s", s)
	}

	key := strings.TrimSpace(choiceParts[0])
	value := strings.TrimSpace(choiceParts[1])

	return key, value, conditions, nil
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

	// Special handling for show commands
	if parts[0] == "show" {
		return ScriptAction{
			Variable: "show",
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
