package aggregate

import "time"

type (
	Event interface {
		ID() string
		Type() string
		OccurredOn() time.Time
	}
	DomainEvent struct {
		id         string
		msgType    string
		occurredAt time.Time
	}
)

func (d DomainEvent) ID() string { return d.id }

func (d DomainEvent) Type() string { return d.msgType }

func (d DomainEvent) OccurredOn() time.Time { return d.occurredAt }