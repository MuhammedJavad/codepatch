package aggregate

type Interface interface {
	Events() []Event
}