package xkafka

import (
	"fmt"
	"github.com/blackRice-Tu/golib"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	asyncProducerMap sync.Map
)

// NewAsyncProducer ...
func NewAsyncProducer(config *Config, version *sarama.KafkaVersion) (*AsyncProducer, error) {
	logger := golib.GetStdLogger()
	if config == nil {
		e := errors.Errorf("config can not be nil")
		logger.Println(e.Error())
		return nil, e
	}
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = golib.GetServerIp()
	if version == nil {
		version = &sarama.V0_10_0_1
	}
	kafkaConfig.Version = *version
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.Retry.Max = 3
	kafkaConfig.Producer.Compression = sarama.CompressionNone
	kafkaConfig.Producer.Flush.Messages = 1000
	kafkaConfig.Producer.Flush.Frequency = time.Second

	// with auth
	if config.SaslUsername != "" {
		kafkaConfig.Net.SASL.User = config.SaslUsername
		kafkaConfig.Net.SASL.Password = config.SaslPassword
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Handshake = true
	}
	brokers := []string{config.Broker}
	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		e := errors.WithMessagef(err, "NewAsyncProducer")
		logger.Println(e.Error())
		return nil, e
	}
	go func() {
		for {
			select {
			case _ = <-producer.Successes():
			case err := <-producer.Errors():
				if err != nil {
					logger.Println(fmt.Sprintf("kafka producer send error: %s", err.Error()))
				}
			}
		}
	}()
	pd := &AsyncProducer{
		id:       config.Id,
		producer: producer,
	}
	return pd, nil
}

// LoadOrNewAsyncProducer ...
func LoadOrNewAsyncProducer(id string, config *Config, version *sarama.KafkaVersion) (*AsyncProducer, error) {
	v, ok := asyncProducerMap.Load(id)
	if ok {
		return v.(*AsyncProducer), nil
	}
	producer, err := NewAsyncProducer(config, version)
	if err != nil {
		return nil, err
	}
	asyncProducerMap.Store(id, producer)
	return producer, nil
}

func (p *AsyncProducer) Send(topic string, msg []byte, key *string) (e error) {
	logger := golib.GetStdLogger()
	if p == nil {
		e = errors.Errorf("producer can not be nil")
		logger.Println(e.Error())
		return
	}
	message := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}
	if key != nil && *key != "" {
		message.Key = sarama.StringEncoder(*key)
	}
	p.producer.Input() <- &message
	return
}

func (p *AsyncProducer) Close() (e error) {
	logger := golib.GetStdLogger()
	if p == nil {
		e = errors.Errorf("producer can not be nil")
		logger.Println(e.Error())
		return
	}
	asyncProducerMap.Delete(p.id)
	return p.producer.Close()
}

func CleanAsyncProducer() {
	asyncProducerMap.Range(func(key, value any) bool {
		pd := value.(*AsyncProducer)
		pd.producer.AsyncClose()
		asyncProducerMap.Delete(key)
		return true
	})
}
