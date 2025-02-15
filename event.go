// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// The canonical Source Code Repository for this Covered Software is:
// https://github.com/gitwyrm/idleyou

package main

type Event struct {
	Name      string
	Done      bool
	Condition func() bool
	Action    func() bool
	Choices   map[string]string
}

// Creates a new Event which is displayed in the eventContainer
//
// If the condition returns true, the event is executed.
// If the action returns true, the event is marked as done.
func NewEvent(name string, condition func() bool, action func() bool, choices map[string]string) Event {
	return Event{
		Name:      name,
		Done:      false,
		Condition: condition,
		Action:    action,
		Choices:   choices,
	}
}

func GetEvents(appstate *AppState) []Event {
	var events []Event

	script := readScript()
	scriptEvents := parseScript(script)
	for _, scriptEvent := range scriptEvents {
		event := scriptEventToEvent(appstate, scriptEvent)
		events = append(events, event)
	}

	otherEvents := []Event{}

	events = append(events, otherEvents...)
	return events
}
