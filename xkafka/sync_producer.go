package xkafka

import (
	"sync"

	"github.com/blackRice-Tu/golib"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	syncProducerMap sync.Map
)

// NewSyncProducer ...
func NewSyncProducer(config *Config, version *sarama.KafkaVersion) (*SyncProducer, error) {
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

	// with auth
	if config.SaslUsername != "" {
		kafkaConfig.Net.SASL.User = config.SaslUsername
		kafkaConfig.Net.SASL.Password = config.SaslPassword
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Handshake = true
	}
	brokers := []string{config.Broker}
	producer, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		e := errors.WithMessagef(err, "NewSyncProducer")
		logger.Println(e.Error())
		return nil, e
	}
	pd := &SyncProducer{
		id:       config.Id,
		producer: producer,
	}
	return pd, nil
}

// LoadOrNewSyncProducer ...
func LoadOrNewSyncProducer(id string, config *Config, version *sarama.KafkaVersion) (*SyncProducer, error) {
	v, ok := syncProducerMap.Load(id)
	if ok {
		return v.(*SyncProducer), nil
	}
	producer, err := NewSyncProducer(config, version)
	if err != nil {
		return nil, err
	}
	syncProducerMap.Store(id, producer)
	return producer, nil
}

// SendMsg ...
func (p *SyncProducer) SendMsg(topic string, msg []byte, key *string) (partition int32, offset int64, e error) {
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
	partition, offset, e = p.producer.SendMessage(&message)
	return
}

// SendMsgList ...
func (p *SyncProducer) SendMsgList(topic string, msgList [][]byte) (e error) {
	logger := golib.GetStdLogger()
	if p == nil {
		e = errors.Errorf("producer can not be nil")
		logger.Println(e.Error())
		return
	}
	messageList := make([]*sarama.ProducerMessage, 0)
	for _, msg := range msgList {
		message := sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msg),
		}
		messageList = append(messageList, &message)
	}
	e = p.producer.SendMessages(messageList)
	return
}

// Close ...
func (p *SyncProducer) Close() (e error) {
	logger := golib.GetStdLogger()
	if p == nil {
		e = errors.Errorf("producer can not be nil")
		logger.Println(e.Error())
		return
	}
	syncProducerMap.Delete(p.id)
	return p.producer.Close()
}

func CleanSyncProducer() {
	logger := golib.GetStdLogger()
	syncProducerMap.Range(func(key, value any) bool {
		pd := value.(*SyncProducer)
		err := pd.producer.Close()
		if err != nil {
			logger.Println(err.Error())
		}
		syncProducerMap.Delete(key)
		return true
	})
}
