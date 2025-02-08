package main

type Event struct {
	Name      string
	Done      bool
	Condition func() bool
	Action    func() bool
}

func NewEvent(name string, condition func() bool, action func() bool) Event {
	return Event{
		Name:      name,
		Done:      false,
		Condition: condition,
		Action:    action,
	}
}

func GetEvents(appstate *AppState) []Event {
	return []Event{
		NewEvent("Promotion to Sales clerk", func() bool { return true }, func() bool {
			appstate.Job.Set("Sales clerk")
			appstate.Messages.Prepend("You got promoted to Sales clerk!")
			return true
		}),
	}
}
