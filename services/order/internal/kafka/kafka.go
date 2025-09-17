package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic: topic,
	})

	return &Producer{writer: writer}
}

func (p *Producer) SendMessage(key string, value string) error {
	msg := kafka.Message{
		Key: []byte(key),
		Value: []byte(value),
	}

	return p.writer.WriteMessages(context.Background(), msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
