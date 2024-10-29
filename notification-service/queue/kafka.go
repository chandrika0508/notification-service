package queue

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Broker    string
	Topic     string
	Partition int
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(config KafkaConfig) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{config.Broker},
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &Producer{writer: writer}
}

func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	msg := kafka.Message{
		Value: message,
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
