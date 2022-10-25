package estream

import (
	"encoding/json"
	"errors"
	"time"
)

//go:generate easyjson

// Meta model
type (
	Meta struct {
		EventName string `json:"event_name"`
		Timestamp int64  `json:"timestamp"`
		UID       string `json:"uid"`
	}

	Payload interface {
		UnmarshalJSON([]byte) error
		PartitionKey() string
	}

	RawMessageHandler func(Meta, json.RawMessage) error

	EventHandler interface {
		RawHandler() RawMessageHandler
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
	UserCreatedPayload struct {
		PublicID  string    `json:"public_id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		Timestamp time.Time `json:"timestamp"`
	}

	UserUpdatedPayload struct {
		PublicID  string    `json:"public_id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		Timestamp time.Time `json:"timestamp"`
	}

	UserDeletedPayload struct {
		PublicID string `json:"public_id"`
	}
)

func (uc *UserCreatedPayload) PartitionKey() string {
	return uc.PublicID
}

func (uc *UserUpdatedPayload) PartitionKey() string {
	return uc.PublicID
}

func (uc *UserDeletedPayload) PartitionKey() string {
	return uc.PublicID
}
