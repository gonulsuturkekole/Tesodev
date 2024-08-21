package consumer

import (
	"context"
	"github.com/labstack/gommon/log"
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
		MaxWait:   time.Second * 1,
		Dialer:    c.dialer,
	})

	c.reader.SetOffset(-1)

}

//config := &kafka.ConfigMap{ "bootstrap.servers": "localhost:9092", "group.id":  "my-group", "auto.offset.reset": "earliest", "enable.auto.commit": true, // Enable auto commit }

func (c *Consumer) Read(callback func(string, error)) {

	// 10 saniyelik bir timeout süresi belirliyoruz
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		// Hata durumunda loglama ve döngüye devam etme
		log.Errorf("Error reading message: %v", err)
		return
	}
	// Mesajın UUID olduğunu varsayarak direkt string olarak ele alıyoruz
	uuid := string(message.Value)

	callback(uuid, nil)

}
