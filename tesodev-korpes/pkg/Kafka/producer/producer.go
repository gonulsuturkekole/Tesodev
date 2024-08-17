package producer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

type Producer struct {
	writer  *kafka.Writer
	topic   string
	brokers []string
}

// NewProducer creates a new Producer instance.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		topic:   topic,
		brokers: brokers,
	}
}

// CreateTopic creates the topic in Kafka if it doesn't exist.
func (p *Producer) CreateTopic() error {
	conn, err := kafka.Dial("tcp", p.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial Kafka broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to dial Kafka controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             p.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("failed to create Kafka topic: %w", err)
	}

	return nil
}

// ProduceMessage creates the topic if it doesn't exist and produces a message to Kafka.
func (p *Producer) ProduceMessage(orderID string) error {
	// Önce topic'i yarat veya varsa atla
	err := p.CreateTopic()
	if err != nil {
		return fmt.Errorf("failed to create topic before producing message: %w", err)
	}

	// Kafka'ya mesaj gönder
	err = p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("OrderID"),
		Value: []byte(orderID),
	})

	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	fmt.Printf("OrderID produced to Kafka: %s\n", orderID)
	return nil
}
