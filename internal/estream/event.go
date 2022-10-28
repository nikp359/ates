package estream

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
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

	TaskCreated    = "task.created"
	TaskUpdated    = "task.updated"
	TaskRegistered = "task.registered"
	TaskCompeted   = "task.completed"
	TasksShuffled  = "task.shuffled"
)

// Map events with topics
var eventTopic = map[string]Topic{
	UserCreated:    TopicUserStreaming,
	UserUpdated:    TopicUserStreaming,
	UserDeleted:    TopicUserStreaming,
	TaskCreated:    TopicTaskStreaming,
	TaskUpdated:    TopicTaskStreaming,
	TaskRegistered: TopicTaskLifecycle,
	TaskCompeted:   TopicTaskLifecycle,
	TasksShuffled:  TopicTaskLifecycle,
}

// ErrUnsupportedEvent for undefined eventName
var ErrUnsupportedEvent = errors.New("unsupported event, see estream.event map for full events list")

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
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserUpdatedPayload struct {
		PublicID  string    `json:"public_id"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserDeletedPayload struct {
		PublicID string `json:"public_id"`
	}

	TaskCreatedPayload struct {
		PublicID       string    `json:"public_id"`
		Title          string    `json:"email"`
		JiraID         string    `json:"role"`
		Description    string    `json:"description"`
		Status         string    `json:"status"`
		AssignedUserID string    `json:"assigned_user_id"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	TaskUpdatedPayload struct {
		PublicID       string    `json:"public_id"`
		Title          string    `json:"email"`
		JiraID         string    `json:"role"`
		Description    string    `json:"description"`
		Status         string    `json:"status"`
		AssignedUserID string    `json:"assigned_user_id"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	TaskRegisteredPayload struct {
		PublicID       string    `json:"public_id"`
		Title          string    `json:"email"`
		JiraID         string    `json:"role"`
		Description    string    `json:"description"`
		Status         string    `json:"status"`
		AssignedUserID string    `json:"assigned_user_id"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	TaskCompletedPayload struct {
		PublicID  string    `json:"public_id"`
		UserID    string    `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	TaskShuffledPayload []TaskAssigned

	TaskAssigned struct {
		UserID string `json:"user_id"`
		TaskID string `json:"task_id"`
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

func (tc *TaskCreatedPayload) PartitionKey() string {
	return tc.PublicID
}

func (tc *TaskUpdatedPayload) PartitionKey() string {
	return tc.PublicID
}

func (tc *TaskRegisteredPayload) PartitionKey() string {
	return tc.PublicID
}

func (tc *TaskCompletedPayload) PartitionKey() string {
	return tc.PublicID
}

func (tc *TaskShuffledPayload) PartitionKey() string {
	return strconv.Itoa(rand.Int())
}
