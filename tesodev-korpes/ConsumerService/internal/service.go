package internal

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"tesodev-korpes/ConsumerService/clientCon"
	"tesodev-korpes/ConsumerService/internal/types"
	_ "tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/pkg"
	"tesodev-korpes/pkg/kafka/consumer"
)

const defaultVatRate = 18.0

type Service struct {
	repo      *Repository
	Consumer  *consumer.Consumer
	conClient *clientCon.ConsumerClient
	vatRate   float64
}

func NewService(repo *Repository, conClient *clientCon.ConsumerClient) *Service {
	// Kafka consumer'ı başlat
	consumer := &consumer.Consumer{Topic: "order-topic"}
	consumer.CreateConnection()

	return &Service{
		repo:      repo,
		Consumer:  consumer,
		conClient: conClient,
		vatRate:   defaultVatRate,
	}
}

func (s *Service) CalculateVat(price float64) float64 {
	return price + (price * s.vatRate / 100)
}
func (s *Service) ConsumeTopic(ctx context.Context, token string) error {
	// Kafka'dan mesajları okuyun
	s.Consumer.Read(func(orderID string, err error) {
		if err != nil {
			log.Errorf("Error reading from Kafka: %v", err)
			return
		}

		log.Infof("Order ID: %s\n", orderID)

		// `conClient` üzerinden sipariş bilgilerini alın
		order, err := s.conClient.GetOrderByID(orderID, token)
		if err != nil {
			log.Errorf("Error getting order by ID: %v", err)
			return
		}
		if order == nil {
			log.Errorf("Order not found for ID: %s", orderID)
			return
		}
		log.Infof("Order Info: %+v", order)

		// Müşteri bilgilerini alın
		customer, err := s.conClient.GetCustomerByID(order.CustomerId, token)
		if err != nil {
			log.Errorf("Error getting customer by ID: %v", err)
			return
		}
		if customer == nil {
			log.Errorf("Customer not found for ID: %s", order.CustomerId)
			return
		}
		log.Infof("Customer Info: %+v", customer)

		priceWithVat := s.CalculateVat(order.Price)
		log.Infof("Price with VAT: %.2f", priceWithVat)

		order.Price = priceWithVat

		// Consumer tipinde bir nesne oluşturun
		consum := &types.Consumer{
			Id:       uuid.New().String(),
			Customer: *customer,
			Order:    *order,
		}

		// Bu nesneyi veritabanına kaydedin
		_, err = s.repo.Create(ctx, consum)
		if err != nil {
			log.Errorf("Error saving consumer to repository: %v", err)
			return
		}

		log.Infof("Consumer saved successfully: %+v", consum)
	})

	return nil
}
