package kk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer  *kafka.Writer
	clients []http.ResponseWriter
	mu      sync.Mutex
}

func NewKafkaWriter(broker, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return p.writer.WriteMessages(ctx, msgs...)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func (kp *KafkaProducer) AddClient(w http.ResponseWriter) {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	kp.clients = append(kp.clients, w)
}

func (kp *KafkaProducer) RemoveClient(w http.ResponseWriter) {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	for i, c := range kp.clients {
		if c == w {
			kp.clients = append(kp.clients[:i], kp.clients[i+1:]...)
			break
		}
	}
}

func (kp *KafkaProducer) Broadcast(msg map[string]interface{}) {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	for _, c := range kp.clients {
		fmt.Fprintf(c, "data: %s\n\n", string(toJSON(msg)))
		c.(http.Flusher).Flush()
	}
}
func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (kp *KafkaProducer) Publish(topic string, message map[string]interface{}) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Failed to marshal message:", err)
		return
	}

	err = kp.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: msgBytes,
		},
	)
	if err != nil {
		fmt.Println("Failed to write message:", err)
		return
	}

	kp.Broadcast(message)
}