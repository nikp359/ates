package estream

const (
	// TopicUser users CUD topic
	TopicUser Topic = "user"
)

// Topic is kafka topic name
type Topic string

// String returns topic name as a string
func (t Topic) String() string {
	return string(t)
}
