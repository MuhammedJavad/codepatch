package aggregate

type Root struct {
	events       []Event
}

func (a *Root) AddDomainNotification(event Event) {
	a.events = append(a.events, event)
}

func (a Root) Events() []Event {
	return a.events
}

