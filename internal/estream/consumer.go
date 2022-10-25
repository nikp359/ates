package estream

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type (
	ConsumerGroup struct {
		client sarama.ConsumerGroup
		topics []string
		ready  chan struct{}
	}
)

func NewConsumerGroup(group string, config Config) (*ConsumerGroup, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = true
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	saramaConfig.ClientID = "estream"
	sarama.Logger = logrus.New()

	client, err := sarama.NewConsumerGroup(config.Addresses, group, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		client: client,
		ready:  make(chan struct{}),
		topics: []string{TopicUserStreaming.String()},
	}, nil
}

func (g *ConsumerGroup) Consume() {
	for {
		if err := g.client.Consume(context.Background(), g.topics, g); err != nil {
			logrus.WithError(err).Error("Error from consumer")
		}
	}
}

func (g *ConsumerGroup) Close() error {
	return g.client.Close()
}

func (g *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	logrus.Infof("Setup sessions")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (g *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (g *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logrus.Infof("Start claim")
	for {
		select {
		case message := <-claim.Messages():
			logrus.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
