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
