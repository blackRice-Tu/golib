package xkafka

import (
	"github.com/Shopify/sarama"
)

type Config struct {
	Id           string `yaml:"id" json:"id"`
	Broker       string `yaml:"broker" json:"broker"`
	SaslUsername string `yaml:"saslUsername" json:"sasl_username"`
	SaslPassword string `yaml:"saslPassword" json:"sasl_password"`
}

type SyncProducer struct {
	id       string
	producer sarama.SyncProducer
}

type AsyncProducer struct {
	id       string
	producer sarama.AsyncProducer
}
