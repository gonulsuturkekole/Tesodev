package internal

import (
	"context"
	"github.com/labstack/gommon/log"
	_ "tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/pkg"
	"tesodev-korpes/pkg/kafka/consumer"
)

type Service struct {
	repo     *Repository
	Consumer *consumer.Consumer // Mesajları string olarak okuyan Kafka Consumer
}

func NewService(repo *Repository) *Service {
	// Kafka consumer'ı başlat
	consumer := &consumer.Consumer{Topic: "order-topic"}
	consumer.CreateConnection()

	return &Service{
		repo:     repo,
		Consumer: consumer,
	}
}

func (s *Service) ConsumeTopic(ctx context.Context) error {
	// Kafka'dan mesajları okuyun
	s.Consumer.Read(func(orderID string, err error) {
		if err != nil {
			log.Errorf("Error reading from Kafka: %v", err)
			return
		}

		log.Infof("Order ID: %s\n", orderID)

		// Burada orderID ile ilgili işlem yapabilirsiniz
		// Örneğin, veritabanına kaydetmek gibi

	})
	return nil
}
