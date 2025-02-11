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

type ScriptCondition struct {
	Variable string
	Operator string
	Value    interface{} // string, float64 or int
}

type ScriptAction struct {
	Variable string
	Operator string
	Value    interface{} // string, float64 or int
}

type ScriptEvent struct {
	Name             string
	ScriptConditions []ScriptCondition
	ScriptActions    []ScriptAction
	Choices          map[string]string
	Return           bool
}

func scriptEventToEvent(state *AppState, scriptEvent ScriptEvent) Event {
	isMultipleChoice := len(scriptEvent.Choices) > 0
	conditions := []func() bool{}
	for _, condition := range scriptEvent.ScriptConditions {
		conditions = append(conditions, scriptConditionToFn(state, condition))
	}
	actions := []func(){}
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

func scriptActionToFn(state *AppState, action ScriptAction, isMultipleChoice bool) func() {
	switch action.Operator {
	case "=":
		return func() {
			state.Set(action.Variable, action.Value)
		}
	case "+=":
		return func() {
			currentValue := state.Get(action.Variable)
			switch currentValue.(type) {
			case int:
				state.Set(action.Variable, currentValue.(int)+action.Value.(int))
			default:
				log.Fatal("Unsupported type for += operation")
			}
		}
	case "-=":
		return func() {
			currentValue := state.Get(action.Variable)
			switch currentValue.(type) {
			case int:
				state.Set(action.Variable, currentValue.(int)-action.Value.(int))
			default:
				log.Fatal("Unsupported type for -= operation")
			}
		}
	case "":
		if action.Variable == "print" {
			return func() {
				if isMultipleChoice {
					// TODO: Display in multiple choice container
				} else {
					state.Messages.Prepend(action.Value.(string))
				}
			}
		}
	}

	return func() {}
}

func scriptConditionToFn(state *AppState, condition ScriptCondition) func() bool {
	valueIsInt := false
	valueIsFloat := false
	switch condition.Value.(type) {
	case int:
		valueIsInt = true
	case float64:
		valueIsFloat = true
	default:
		valueIsInt = false
		valueIsFloat = false
	}

	//variableIsInt := false
	variableIsFloat := false
	//variableIsString := false

	switch state.Get(condition.Variable).(type) {
	case int:
		//variableIsInt = true
	case float64:
		variableIsFloat = true
	case string:
		//variableIsString = true
	default:
		//variableIsInt = false
		variableIsFloat = false
		//variableIsString = false
	}

	getVariableValue := func() (int, float64, string) {
		var variableInt int
		var variableFloat float64
		var variableString string

		variable := state.Get(condition.Variable)
		switch variable.(type) {
		case int:
			variableInt = variable.(int)
		case float64:
			variableFloat = variable.(float64)
		case string:
			variableString = variable.(string)
		default:
			return 0, 0.0, ""
		}

		return variableInt, variableFloat, variableString
	}

	switch condition.Operator {
	case "<":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat < condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt < condition.Value.(int)
		}
	case ">":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat > condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt > condition.Value.(int)
		}
	case "<=":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat <= condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt <= condition.Value.(int)
		}
	case ">=":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat >= condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt >= condition.Value.(int)
		}
	case "==":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat == condition.Value.(float64)
			}
		}
		if !valueIsInt {
			_, _, variableString := getVariableValue()
			return func() bool { return variableString == condition.Value.(string) }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt == condition.Value.(int)
		}
	case "!=":
		if valueIsFloat && variableIsFloat {
			return func() bool {
				_, variableFloat, _ := getVariableValue()
				return variableFloat != condition.Value.(float64)
			}
		}
		if !valueIsInt {
			_, _, variableString := getVariableValue()
			return func() bool { return variableString != condition.Value.(string) }
		}
		return func() bool {
			variableInt, _, _ := getVariableValue()
			return variableInt != condition.Value.(int)
		}
	default:
		panic(fmt.Sprintf("Unknown operator: %s", condition.Operator))
	}
}

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
				Return := line[2:]
				switch Return {
				case "true":
					currentEvent.Return = true
				case "false":
					currentEvent.Return = false
				default:
					log.Fatal("Error parsing ScriptEvent Return")
				}
			}
		}
	}

	// Append last event
	if currentEvent != nil {
		events = append(events, *currentEvent)
	}

	return events
}

// Parses a choice line with the format:
// key -> value
func parseChoice(s string) (string, string) {
	parts := strings.Split(s, "->")
	if len(parts) != 2 {
		panic(fmt.Sprintf("Invalid choice syntax: %s", s))
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	return key, value
}

func parseCondition(line string) ScriptCondition {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) != 3 {
		panic(fmt.Sprintf("Invalid condition syntax: %s", line))
	}

	value, err := strconv.Atoi(parts[2])
	if err != nil {
		return ScriptCondition{
			Variable: parts[0],
			Operator: parts[1],
			Value:    parts[2],
		}
	}

	return ScriptCondition{
		Variable: parts[0],
		Operator: parts[1],
		Value:    value,
	}
}

func parseAction(line string) ScriptAction {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) < 2 {
		panic(fmt.Sprintf("Invalid action syntax: %s", line))
	}

	if parts[0] == "print" { // Special case for print statements
		return ScriptAction{
			Variable: "print",
			Operator: "",
			Value:    strings.Join(parts[1:], " "),
		}
	}

	value, err := strconv.Atoi(parts[2])
	if err != nil {
		return ScriptAction{
			Variable: parts[0],
			Operator: parts[1],
			Value:    strings.Join(parts[2:], " "),
		}
	}

	return ScriptAction{
		Variable: parts[0],
		Operator: parts[1],
		Value:    value,
	}
}
