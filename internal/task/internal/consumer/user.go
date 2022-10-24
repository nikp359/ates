package consumer

import "github.com/nikp359/ates/internal/task/internal/repository"

type (
	UserConsumer struct {
		userRepository *repository.UserRepository
	}
)

func NewUserConsumer(userRepository *repository.UserRepository) *UserConsumer {
	return &UserConsumer{
		userRepository: userRepository,
	}
}

func (c *UserConsumer) Start() {

}
