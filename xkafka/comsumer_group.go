package xkafka

import (
	"context"
	"sync"

	"github.com/blackRice-Tu/golib"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type (
	ConsumerGroupSession = sarama.ConsumerGroupSession
	ConsumerGroupClaim   = sarama.ConsumerGroupClaim
)

type ConsumerHandler struct {
	setup        func(ConsumerGroupSession) error
	cleanup      func(ConsumerGroupSession) error
	consumeClaim func(ConsumerGroupSession, ConsumerGroupClaim) error
}

func (t ConsumerHandler) Setup(session ConsumerGroupSession) error {
	if t.setup != nil {
		return t.setup(session)
	}
	return nil
}

func (t ConsumerHandler) Cleanup(session ConsumerGroupSession) error {
	if t.cleanup != nil {
		return t.cleanup(session)
	}
	return nil
}

func (t ConsumerHandler) ConsumeClaim(session ConsumerGroupSession, claim ConsumerGroupClaim) error {
	if t.consumeClaim != nil {
		return t.consumeClaim(session, claim)
	}
	return nil
}

type ConsumerGroup struct {
	consumerGroup sarama.ConsumerGroup
	groupId       string
	handler       ConsumerHandler
}

func NewConsumerGroup(config *Config, groupId string, version *sarama.KafkaVersion) (*ConsumerGroup, error) {
	kafkaConfig := sarama.NewConfig()

	if version == nil {
		version = &sarama.V0_10_0_1
	}
	kafkaConfig.Version = *version
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	brokers := []string{config.Broker}
	group, err := sarama.NewConsumerGroup(brokers, groupId, kafkaConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "NewConsumerGroup")
	}

	consumerGroup := &ConsumerGroup{
		consumerGroup: group,
		groupId:       groupId,
		handler:       ConsumerHandler{},
	}
	return consumerGroup, nil
}

func (t *ConsumerGroup) Setup(ctx context.Context, f func(ConsumerGroupSession) error) {
	t.handler.setup = f
}

func (t *ConsumerGroup) Cleanup(ctx context.Context, f func(ConsumerGroupSession) error) {
	t.handler.cleanup = f
}

func (t *ConsumerGroup) ConsumeClaim(ctx context.Context, f func(ConsumerGroupSession, ConsumerGroupClaim) error) {
	t.handler.consumeClaim = f
}

func (t *ConsumerGroup) Consume(ctx context.Context, topics []string) error {
	logger := golib.GetStdLogger()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := t.consumerGroup.Consume(ctx, topics, t.handler); err != nil {
				if errors.Is(err, sarama.ErrClosedClient) || errors.Is(err, sarama.ErrClosedConsumerGroup) {
					// quit
					logger.Println("kafka consumer quit; ", t.groupId)
					return
				}
				logger.Println("kafka consumer error; ", err.Error())
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logger.Println("kafka consumer context was cancelled")
				return
			}
		}
	}()
	wg.Wait()
	return nil
}
