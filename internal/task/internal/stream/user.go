package stream

import (
	"github.com/nikp359/ates/internal/estream"
	"github.com/nikp359/ates/internal/task/internal/repository"
)

type (
	UserStream struct {
		userRepository *repository.UserRepository
		consumer       *estream.ConsumerGroup
	}
)

func NewUserStream(userRepository *repository.UserRepository, consumerConfig estream.Config) (*UserStream, error) {
	consumer, err := estream.NewConsumerGroup("task:user:stream", consumerConfig)
	if err != nil {
		return nil, err
	}

	return &UserStream{
		userRepository: userRepository,
		consumer:       consumer,
	}, nil
}

func (c *UserStream) Start() {
	go c.consumer.Consume()
}
