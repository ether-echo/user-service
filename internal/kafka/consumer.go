package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/ether-echo/user-service/pkg/logger"

	"github.com/IBM/sarama"
)

var (
	log = logger.Logger().Named("kafka").Sugar()
)

type Consumer struct {
	group   sarama.ConsumerGroup
	config  *sarama.Config
	topics  []string
	handler sarama.ConsumerGroupHandler
}

func NewConsumer(brokers, topics []string, groupID string, handler sarama.ConsumerGroupHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group: %w", err)
	}

	return &Consumer{
		group:   group,
		config:  config,
		topics:  topics,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info("Consumer shutting down")
				return
			default:
				err := c.group.Consume(ctx, c.topics, c.handler)
				if err != nil {
					log.Errorf("Error from consumer: %v", err)
				}
			}
		}
	}()

	wg.Wait()
	return nil
}

func (c *Consumer) Stop() error {
	return c.group.Close()
}

func (c *Consumer) Subscribe(topic string) error {
	c.topics = append(c.topics, topic)
	return nil
}
