package main

import (
	"testing"
)

// -------------------------
// Arithmetic Tests
// -------------------------

func TestGameVariable_Add(t *testing.T) {
	a := NewGameVariable("var", 10)
	b := NewGameVariable("var", 5)
	result := a.Add(b)
	if result.Int() != 15 {
		t.Errorf("Expected 15, got %d", result.Int())
	}

	f1 := NewGameVariable("var", 2.5)
	f2 := NewGameVariable("var", 2.5)
	result = f1.Add(f2)
	if result.Float64() != 5.0 {
		t.Errorf("Expected 5.0, got %f", result.Float64())
	}

	s1 := NewGameVariable("var", "Hello ")
	s2 := NewGameVariable("var", "World!")
	result = s1.Add(s2)
	if result.String() != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got '%s'", result.String())
	}
}

func TestGameVariable_Subtract(t *testing.T) {
	a := NewGameVariable("var", 10)
	b := NewGameVariable("var", 5)
	result := a.Subtract(b)
	if result.Int() != 5 {
		t.Errorf("Expected 5, got %d", result.Int())
	}

	f1 := NewGameVariable("var", 5.5)
	f2 := NewGameVariable("var", 2.5)
	result = f1.Subtract(f2)
	if result.Float64() != 3.0 {
		t.Errorf("Expected 3.0, got %f", result.Float64())
	}
}

func TestGameVariable_Multiply(t *testing.T) {
	a := NewGameVariable("var", 4)
	b := NewGameVariable("var", 3)
	result := a.Multiply(b)
	if result.Int() != 12 {
		t.Errorf("Expected 12, got %d", result.Int())
	}

	f1 := NewGameVariable("var", 1.5)
	f2 := NewGameVariable("var", 2.0)
	result = f1.Multiply(f2)
	if result.Float64() != 3.0 {
		t.Errorf("Expected 3.0, got %f", result.Float64())
	}
}

func TestGameVariable_Divide(t *testing.T) {
	a := NewGameVariable("var", 10)
	b := NewGameVariable("var", 2)
	result := a.Divide(b)
	if result.Int() != 5 {
		t.Errorf("Expected 5, got %d", result.Int())
	}

	f1 := NewGameVariable("var", 7.5)
	f2 := NewGameVariable("var", 2.5)
	result = f1.Divide(f2)
	if result.Float64() != 3.0 {
		t.Errorf("Expected 3.0, got %f", result.Float64())
	}

	// Division by zero should not crash
	bZero := NewGameVariable("var", 0)
	result = a.Divide(bZero)
	if result.Int() != 10 {
		t.Errorf("Expected original value 10, got %d", result.Int())
	}
}

// -------------------------
// Comparison Tests
// -------------------------

func TestGameVariable_Compare(t *testing.T) {
	a := NewGameVariable("var", 10)
	b := NewGameVariable("var", 5)

	if !a.Compare(b, ">") {
		t.Errorf("Expected 10 > 5 to be true")
	}
	if a.Compare(b, "<") {
		t.Errorf("Expected 10 < 5 to be false")
	}
	if !a.Compare(a, "==") {
		t.Errorf("Expected 10 == 10 to be true")
	}
	if a.Compare(b, "==") {
		t.Errorf("Expected 10 == 5 to be false")
	}
	if a.Compare(b, "<=") {
		t.Errorf("Expected 10 <= 5 to be false")
	}
	if !a.Compare(b, ">=") {
		t.Errorf("Expected 10 >= 5 to be true")
	}
}

// -------------------------
// Type Conversion Tests
// -------------------------

func TestGameVariable_String(t *testing.T) {
	s := NewGameVariable("var", "hello")
	if s.String() != "hello" {
		t.Errorf("Expected 'hello', got '%s'", s.String())
	}

	i := NewGameVariable("var", 42)
	if i.String() != "42" {
		t.Errorf("Expected '42', got '%s'", i.String())
	}

	f := NewGameVariable("var", 3.14)
	if f.String() != "3.140000" { // fmt.Sprintf("%f") behavior
		t.Errorf("Expected '3.140000', got '%s'", f.String())
	}

	b := NewGameVariable("var", true)
	if b.String() != "true" {
		t.Errorf("Expected 'true', got '%s'", b.String())
	}
}

func TestGameVariable_Int(t *testing.T) {
	i := NewGameVariable("var", 42)
	if i.Int() != 42 {
		t.Errorf("Expected 42, got %d", i.Int())
	}

	f := NewGameVariable("var", 3.14)
	if f.Int() != 0 {
		t.Errorf("Expected 0 (invalid conversion), got %d", f.Int())
	}
}

func TestGameVariable_Float64(t *testing.T) {
	f := NewGameVariable("var", 3.14)
	if f.Float64() != 3.14 {
		t.Errorf("Expected 3.14, got %f", f.Float64())
	}

	i := NewGameVariable("var", 42)
	if i.Float64() != 0.0 {
		t.Errorf("Expected 0.0 (invalid conversion), got %f", i.Float64())
	}
}

func TestGameVariable_Bool(t *testing.T) {
	bTrue := NewGameVariable("var", true)
	if !bTrue.Bool() {
		t.Errorf("Expected true, got false")
	}

	bFalse := NewGameVariable("var", false)
	if bFalse.Bool() {
		t.Errorf("Expected false, got true")
	}

	i := NewGameVariable("var", 42)
	if i.Bool() {
		t.Errorf("Expected false (invalid conversion), got true")
	}
}
