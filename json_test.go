package main

import (
	"testing"
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
}
