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

	SyncProducer struct {
		producer sarama.SyncProducer
	}

	AsyncProducer struct {
		producer sarama.AsyncProducer
	}
)

func NewSyncProducer(config Config) (*SyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	saramaCfg.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Addresses, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		producer: producer,
	}, nil
}

func (sp *SyncProducer) Send(eventName string, payload json.Unmarshaler) error {
	msg, err := getMessage(eventName, payload)
	if err != nil {
		return err
	}

	_, _, err = sp.producer.SendMessage(msg)

	return err
}

func NewAsyncProducer(config Config) (*AsyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	saramaCfg.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	saramaCfg.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(config.Addresses, saramaCfg)
	if err != nil {
		return nil, err
	}

	go func() {
		for err = range producer.Errors() {
			logrus.WithError(err).Error("producer write error")
		}
	}()

	return &AsyncProducer{
		producer: producer,
	}, nil
}

func (ap *AsyncProducer) Send(eventName string, payload json.Unmarshaler) error {
	msg, err := getMessage(eventName, payload)
	if err != nil {
		return err
	}

	ap.producer.Input() <- msg

	return nil
}

func getMessage(eventName string, payload json.Unmarshaler) (*sarama.ProducerMessage, error) {
	t, ok := EventTopic(eventName)
	if !ok {
		return nil, ErrUnsupportedEvent
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
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: t.String(),
		Value: sarama.ByteEncoder(msgBody),
	}, nil
}

func randomUUID() string {
	uid, err := uuid.GenerateUUID()
	if err != nil {
		logrus.Errorf("Generate UUID: %s", err)
		return fmt.Sprintf("fail:%d", time.Now().UnixNano())
	}
	return uid
}
