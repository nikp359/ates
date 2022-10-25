package estream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type (
	ConsumerGroup struct {
		client   sarama.ConsumerGroup
		topics   []string
		handlers map[string]map[string][]RawMessageHandler
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
		client:   client,
		topics:   make([]string, 0),
		handlers: make(map[string]map[string][]RawMessageHandler),
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

// AddHandler add handler for event
func (g *ConsumerGroup) AddHandler(eventName string, h EventHandler) error {
	topic, ok := EventTopic(eventName)
	if !ok {
		return ErrUnsupportedEvent
	}

	if _, ok = g.handlers[topic.String()]; !ok {
		g.handlers[topic.String()] = make(map[string][]RawMessageHandler)
	}

	g.handlers[topic.String()][eventName] = append(g.handlers[topic.String()][eventName], h.RawHandler())
	g.topics = append(g.topics, topic.String())
	return nil
}

func (g *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (g *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (g *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	eventHandlers, ok := g.handlers[claim.Topic()]
	if !ok {
		return fmt.Errorf("missing handler for topic: %s", claim.Topic())
	}

	for {
		select {
		case message := <-claim.Messages():
			var event consumerEvent
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				if len(message.Value) == 0 {
					logrus.WithError(err).Errorf("Empty message error, msg: %+v", message)
					continue
				}
				if !json.Valid(message.Value) {
					logrus.WithError(err).Errorf("Validate json error, msg: %+v msg.Value: %s", message, string(message.Value))
					continue
				}

				return fmt.Errorf("msg: %+v | value: %s | unmarshal estream event: %w", message, string(message.Value), err)
			}

			handlers, handlersFound := eventHandlers[event.EventName]
			if !handlersFound {
				session.MarkMessage(message, "")
				continue
			}

			for _, h := range handlers {
				err = h(event.Meta, event.Payload)
				if err != nil {
					logrus.WithError(err).Errorf("Handle message. Event: %s", event.EventName)
					continue
				}
			}
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
