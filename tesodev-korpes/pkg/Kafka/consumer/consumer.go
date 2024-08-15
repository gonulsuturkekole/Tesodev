package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Consumer struct {
	reader *kafka.Reader
	dialer *kafka.Dialer
	Topic  string
}

func (c *Consumer) CreateConnection() {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     c.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Second * 10,
		Dialer:    c.dialer,
	})

	c.reader.SetOffset(0)
}

func (c *Consumer) Read(callback func(string, error)) {
	for {
		ctx := context.Background()
		message, err := c.reader.ReadMessage(ctx)

		if err != nil {
			callback("", err)
			return
		}

		// Mesajın UUID olduğunu varsayarak direkt string olarak ele alıyoruz
		uuid := string(message.Value)

		callback(uuid, nil)
	}
}
