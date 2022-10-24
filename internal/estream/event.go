package estream

import (
	"encoding/json"
	"errors"
)

//go:generate easyjson

// Meta model
type (
	Meta struct {
		EventName string `json:"event_name"`
		Timestamp int64  `json:"timestamp"`
		UID       string `json:"uid"`
	}

	// Payload interface must implement all events payload models
	Payload interface {
		UnmarshalJSON([]byte) error
	}
)

const (
	UserCreated = "user.create"
	UserUpdated = "user.updated"
	UserDeleted = "user.deleted"
)

// Map events with topics
var eventTopic = map[string]Topic{
	UserCreated: TopicUserStreaming,
	UserUpdated: TopicUserStreaming,
	UserDeleted: TopicUserStreaming,
}

// ErrUnsupportedEvent for undefined eventName
var ErrUnsupportedEvent = errors.New("unsupported event, see estream.eventTopic map for full events list")

// EventTopic by event name
func EventTopic(name string) (Topic, bool) {
	t, ok := eventTopic[name]
	return t, ok
}

type consumerEvent struct {
	Meta
	Payload json.RawMessage `json:"payload"`
}

type producerEvent struct {
	Meta
	Payload json.Unmarshaler `json:"payload"`
}

// easyjson:json
type (
	EventUserCreated struct {
		PublicID string `json:"public_id"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	EventUserUpdated struct {
		Role string `json:"role"`
	}

	EventUserDeleted struct {
		PublicID string `json:"public_id"`
	}
)
