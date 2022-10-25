package stream

import (
	"github.com/nikp359/ates/internal/estream"
	"github.com/nikp359/ates/internal/task/internal/model"
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

	us := &UserStream{
		userRepository: userRepository,
		consumer:       consumer,
	}

	err = consumer.AddHandler(estream.UserCreated, estream.UserCreatedHandler(us.createUser))
	if err != nil {
		return nil, err
	}

	err = consumer.AddHandler(estream.UserUpdated, estream.UserUpdatedHandler(us.updateUser))
	if err != nil {
		return nil, err
	}

	err = consumer.AddHandler(estream.UserDeleted, estream.UserDeletedHandler(us.deleteUser))
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (c *UserStream) Start() {
	go c.consumer.Consume()
}

func (c *UserStream) createUser(_ estream.Meta, event estream.UserCreatedPayload) error {
	return c.userRepository.Add(&model.User{
		PublicID:  event.PublicID,
		Email:     event.Email,
		Role:      event.Role,
		UpdatedAt: event.UpdatedAt,
	})
}

func (c *UserStream) updateUser(_ estream.Meta, event estream.UserUpdatedPayload) error {
	return c.userRepository.Update(&model.User{
		PublicID:  event.PublicID,
		Email:     event.Email,
		Role:      event.Role,
		UpdatedAt: event.UpdatedAt,
	})
}

func (c *UserStream) deleteUser(_ estream.Meta, event estream.UserDeletedPayload) error {
	return c.userRepository.Delete(event.PublicID)
}
