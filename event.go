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

	otherEvents := []Event{
		NewEvent(
			"Promotion to Manager",
			func() bool {
				workXP, err := appstate.WorkXP.Get()
				if err != nil {
					return false
				}
				return workXP == 1000
			},
			func() bool {
				appstate.Job.Set("Manager")
				appstate.Messages.Prepend("Event: You got promoted to Manager!")
				return true
			},
			nil,
		),
		NewEvent(
			"Test Choice Event",
			func() bool {
				return appstate.Get("workxp").(int) == 300
			},
			func() bool {
				return true
			},
			map[string]string{
				"Choice 1": "Chose 1",
				"Choice 2": "Chose 2",
			},
		),
		NewEvent(
			"Chose 1",
			func() bool {
				return false
			},
			func() bool {
				appstate.Messages.Prepend("Event: You chose 1.")
				return true
			},
			nil,
		),
		NewEvent(
			"Chose 2",
			func() bool {
				return false
			},
			func() bool {
				appstate.Messages.Prepend("Event: You chose 2.")
				return true
			},
			nil,
		),
	}

	events = append(events, otherEvents...)
	return events
}
