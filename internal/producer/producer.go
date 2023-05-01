package producer

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type ProducerProcessing struct {
	kafkaProducer *kafka.Conn
	topic         string
}

func NewProducerProcessing(kafkaProducer *kafka.Conn, topic string) ProducerProcessing {
	return ProducerProcessing{
		kafkaProducer: kafkaProducer,
		topic:         topic,
	}
}

func (p ProducerProcessing) WriteMessages(message interface{}, topic ...string) error {
	messageBytes, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to marshal message: %v", err)
	}

	if _, err := p.kafkaProducer.WriteMessages(kafka.Message{
		Topic: p.topic,
		Value: messageBytes,
	}); err != nil {
		return fmt.Errorf("unable to write message into queue: %v", err)
	}

	return nil
}
