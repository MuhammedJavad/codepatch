package box

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID          int64      `json:"id"`
	Data        []byte     `json:"data"`
	ReceivedAt  time.Time  `json:"received_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	Source      string     `json:"source"`
	Event       string     `json:"event"`
	Error       *string    `json:"error"`
	RetryCount  int        `json:"retry_count"`
	RetryMax    int        `json:"retry_max"`
}

func (m *Message) confirm(err error) {
	if err != nil {
		errString := err.Error()
		if len(errString) > 300 {
			errString = errString[:300]
		}
		m.Error = &errString
		m.RetryCount++
		if m.RetryMax == -1 {
			return
		}
		if m.RetryCount < m.RetryMax {
			return
		}
	} else {
		m.Error = nil
	}
	now := time.Now()
	m.ProcessedAt = &now
}

func NewMessage[T interface{}](data T, source, event string, retryMax int) (Message, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return Message{}, err
	}

	m := Message{
		Data:       payload,
		Source:     source,
		Event:      event,
		RetryMax:   retryMax,
		ReceivedAt: time.Now(),
	}
	return m, nil
}
