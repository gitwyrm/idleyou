package main

type Event struct {
	Name      string
	Done      bool
	Condition func() bool
	Action    func() bool
}

// Creates a new Event which is displayed in the eventContainer
//
// If the condition returns true, the event is executed.
// If the action returns true, the event is marked as done.
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
		NewEvent(
			"Promotion to Sales clerk",
			func() bool {
				workXP, err := appstate.WorkXP.Get()
				if err != nil {
					return false
				}
				return workXP == 200
			}, func() bool {
				appstate.Job.Set("Sales clerk")
				appstate.Messages.Prepend("You got promoted to Sales clerk!")
				return true
			}),
	}
}
