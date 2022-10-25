package stream

import (
	"github.com/nikp359/ates/internal/estream"
	"github.com/nikp359/ates/internal/task/internal/repository"
	"github.com/sirupsen/logrus"
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

	us := &UserStream{
		userRepository: userRepository,
		consumer:       consumer,
	}

	err = consumer.AddHandler(estream.UserCreated, estream.UserCreatedHandler(us.createUser))
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (c *UserStream) Start() {
	go c.consumer.Consume()
}

func (c *UserStream) createUser(meta estream.Meta, user estream.UserCreatedPayload) error {
	logrus.Infof("Meta: %+v", meta)
	logrus.Infof("Estream: %+v", user)

	return nil
}
