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
	Return           bool
}

func scriptActionToFn(state *AppState, action ScriptAction) func() {
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
		return func() bool { return false }
	}

	switch condition.Operator {
	case "<":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat < condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			return variableInt < condition.Value.(int)
		}
	case ">":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat > condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			return variableInt > condition.Value.(int)
		}
	case "<=":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat <= condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			return variableInt <= condition.Value.(int)
		}
	case ">=":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat >= condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return false }
		}
		return func() bool {
			return variableInt >= condition.Value.(int)
		}
	case "==":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat == condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return variableString == condition.Value.(string) }
		}
		return func() bool {
			return variableInt == condition.Value.(int)
		}
	case "!=":
		if valueIsFloat && variableFloat != 0.0 {
			return func() bool {
				return variableFloat != condition.Value.(float64)
			}
		}
		if !valueIsInt {
			return func() bool { return variableString != condition.Value.(string) }
		}
		return func() bool {
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
			Value:    parts[2],
		}
	}

	return ScriptAction{
		Variable: parts[0],
		Operator: parts[1],
		Value:    value,
	}
}
