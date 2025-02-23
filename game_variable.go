package main

import (
	"fmt"
	"strconv"
)

type GameVariable struct {
	// name of the variable in AppState
	name string
	// string, int, float64, bool
	value interface{}
}

func NewGameVariable(name string, value interface{}) GameVariable {
	return GameVariable{name, value}
}

func NewGameVariableFromString(name, value string) GameVariable {
	// int
	if i, err := strconv.Atoi(value); err == nil {
		return NewGameVariable(name, i)
	}
	// float64
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return NewGameVariable(name, f)
	}
	// bool
	if b, err := strconv.ParseBool(value); err == nil {
		return NewGameVariable(name, b)
	}
	// string
	if s := value; s != "" {
		return NewGameVariable(name, s)
	}
	return NewGameVariable(name, value)
}

// ------------------------------
// Mathematical operations
// ------------------------------

// Add two GameVariables
func (gv GameVariable) Add(other GameVariable) GameVariable {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok {
			return GameVariable{gv.name, v + ov}
		}
	case float64:
		if ov, ok := other.value.(float64); ok {
			return GameVariable{gv.name, v + ov}
		}
	case string:
		if ov, ok := other.value.(string); ok {
			return GameVariable{gv.name, v + ov}
		}
	case bool:
		if ov, ok := other.value.(bool); ok {
			return GameVariable{gv.name, v || ov}
		}
	}
	return gv
}

// Subtract two GameVariables
func (gv GameVariable) Subtract(other GameVariable) GameVariable {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok {
			return GameVariable{gv.name, v - ov}
		}
	case float64:
		if ov, ok := other.value.(float64); ok {
			return GameVariable{gv.name, v - ov}
		}
	}
	return gv
}

// Multiply two GameVariables
func (gv GameVariable) Multiply(other GameVariable) GameVariable {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok {
			return GameVariable{gv.name, v * ov}
		}
	case float64:
		if ov, ok := other.value.(float64); ok {
			return GameVariable{gv.name, v * ov}
		}
	}
	return gv
}

// Divide two GameVariables
func (gv GameVariable) Divide(other GameVariable) GameVariable {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok && ov != 0 {
			return GameVariable{gv.name, v / ov}
		}
	case float64:
		if ov, ok := other.value.(float64); ok && ov != 0 {
			return GameVariable{gv.name, v / ov}
		}
	}
	return gv
}

// ------------------------------
// Comparison
// ------------------------------

func (gv GameVariable) SameType(other GameVariable) bool {
	switch gv.value.(type) {
	case int:
		if _, ok := other.value.(int); ok {
			return true
		}
	case float64:
		if _, ok := other.value.(float64); ok {
			return true
		}
	case string:
		if _, ok := other.value.(string); ok {
			return true
		}
	case bool:
		if _, ok := other.value.(bool); ok {
			return true
		}
	}
	return false
}

func (gv GameVariable) LesserThan(other GameVariable) bool {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok {
			return v < ov
		}
	case float64:
		if ov, ok := other.value.(float64); ok {
			return v < ov
		}
	}
	return false
}

func (gv GameVariable) GreaterThan(other GameVariable) bool {
	switch v := gv.value.(type) {
	case int:
		if ov, ok := other.value.(int); ok {
			return v > ov
		}
	case float64:
		if ov, ok := other.value.(float64); ok {
			return v > ov
		}
	}
	return false
}

// Compare two GameVariables
func (gv GameVariable) Compare(other GameVariable, operator string) bool {
	switch operator {
	case "==":
		return gv.value == other.value
	case "!=":
		return gv.value != other.value
	case "<":
		return gv.LesserThan(other)
	case ">":
		return gv.GreaterThan(other)
	case "<=":
		return gv.LesserThan(other) || gv.value == other.value
	case ">=":
		return gv.GreaterThan(other) || gv.value == other.value
	default:
		return false
	}
}

// ------------------------------
// Casting
// ------------------------------

func (gv GameVariable) String() string {
	switch v := gv.value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%f", v)
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return "unknown"
	}
}

func (gv GameVariable) Bool() bool {
	if v, ok := gv.value.(bool); ok {
		return v
	}
	return false
}

func (gv GameVariable) Int() int {
	if v, ok := gv.value.(int); ok {
		return v
	}
	return 0
}

func (gv GameVariable) Float64() float64 {
	if v, ok := gv.value.(float64); ok {
		return v
	}
	return 0.0
}

// ------------------------------
// AppState
// ------------------------------

func (gv GameVariable) UpdateAppState(state *AppState) {
	state.Set(gv.name, gv.value)
}
