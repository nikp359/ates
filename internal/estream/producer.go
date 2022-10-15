package estream

import (
	"github.com/Shopify/sarama"
)

type (
	Config struct {
		Addresses []string `yaml:"addresses"`
	}
)

func NewSyncProducer(config Config) (sarama.SyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	saramaCfg.Producer.Compression = sarama.CompressionSnappy // Compress messages
	saramaCfg.Producer.Return.Successes = true

	return sarama.NewSyncProducer(config.Addresses, saramaCfg)
}
