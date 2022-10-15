package estream

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Addresses []string `yaml:"addresses"`
	}

	Producer struct {
		syncProducer sarama.SyncProducer
	}
)

func NewSyncProducer(config Config) (*Producer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	saramaCfg.Producer.Compression = sarama.CompressionSnappy // Compress messages
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Addresses, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		syncProducer: producer,
	}, nil
}

func (p *Producer) SendSync(eventName string, payload json.Unmarshaler) error {
	t, ok := EventTopic(eventName)
	if !ok {
		return ErrUnsupportedEvent
	}

	msg := &producerEvent{
		Meta: Meta{
			EventName: eventName,
			Timestamp: time.Now().Unix(),
			UID:       randomUUID(),
		},
		Payload: payload,
	}

	msgBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: t.String(),
		Value: sarama.ByteEncoder(msgBody),
	})

	return err
}

func randomUUID() string {
	uid, err := uuid.GenerateUUID()
	if err != nil {
		logrus.Errorf("Generate UUID: %s", err)
		return fmt.Sprintf("fail:%d", time.Now().UnixNano())
	}
	return uid
}
