package estream

const (
	TopicUserAuth      Topic = "user.auth"
	TopicUserStreaming Topic = "user.streaming"
)

// Topic is kafka topic name
type Topic string

// String returns topic name as a string
func (t Topic) String() string {
	return string(t)
}
