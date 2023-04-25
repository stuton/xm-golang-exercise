package producer

import (
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/stuton/xm-golang-exercise/internal/model"
)

type EventMessage struct {
	EventType EventType
	Company   model.Company
}

type ProducerProcessing struct {
	kafkaProducer *kafka.Conn
}

func NewProducerProcessing(kafkaProducer *kafka.Conn) ProducerProcessing {
	return ProducerProcessing{
		kafkaProducer: kafkaProducer,
	}
}

func (p ProducerProcessing) WriteMessages(topic string, message interface{}) error {
	messageBytes, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("unable to marshal message: %v", err)
	}

	if _, err := p.kafkaProducer.WriteMessages(kafka.Message{
		Topic: topic,
		Value: messageBytes,
	}); err != nil {
		return fmt.Errorf("unable to write message into queue: %v", err)
	}

	return nil
}
