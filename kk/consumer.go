package kk

import "github.com/segmentio/kafka-go"

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker, groupID, topic string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID, // Consumer group ID
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) Reader() *kafka.Reader {
	return c.reader
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
