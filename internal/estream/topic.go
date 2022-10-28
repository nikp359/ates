package estream

const (
	TopicUserStreaming Topic = "user.streaming" // for CUD events
	TopicTaskStreaming Topic = "task.streaming" // for CUD events
	TopicTaskLifecycle Topic = "task.lifecycle" // for business events
)

// Topic is kafka topic name
type Topic string

// String returns topic name as a string
func (t Topic) String() string {
	return string(t)
}
