package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"tesodev-korpes/ConsumerService/clientCon"
	"tesodev-korpes/ConsumerService/internal/types"
	"tesodev-korpes/pkg/Kafka/consumer"
)

type Service struct {
	repo          *Repository
	kafkaConsumer *consumer.Consumer
	conClient     *clientCon.ConsumerClient
}

func NewService(repo *Repository, conClient *clientCon.ConsumerClient, kafkaConsumer *consumer.Consumer, brokers []string, topic string) *Service {

	kafkaConsumer.Topic = topic
	kafkaConsumer.CreateConnection(brokers)

	return &Service{
		repo:          repo,
		conClient:     conClient,
		kafkaConsumer: kafkaConsumer,
	}
}

func (s *Service) ProcessMessage(ctx context.Context, msg string) error {
	// Mesajı işleyin ve Order Service'e istek gönderin
	err := s.sendRequest(ctx, msg)
	if err != nil {
		fmt.Printf("Error sending order request: %v\n", err)
		return err
	}
	return nil
}

func (s *Service) sendRequest(ctx context.Context, msg string) error {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjZiNmEzNWUxLTgyNjYtNDViMi05YTc2LTMxOGE3YTVjMzE0NiIsImZpcnN0X25hbWUiOiJBeXNlIiwibGFzdF9uYW1lIjoiQ2FuIiwiZXhwIjoxNzI0Mzk5NDUwfQ.EkjZOCooaTbGnuY6zVbuJf9mBxc1VjBAg_MKG5Xr2Mo"
	order, err := s.conClient.GetOrderByID(msg, token)
	if err != nil {
		log.Errorf("Error getting order by ID: %v", err)
		return nil
	}
	if order == nil {
		log.Errorf("Order not found for ID: %s", msg)
		return nil
	}
	log.Infof("Order Info: %+v", order)

	// Müşteri bilgilerini alın
	customer, err := s.conClient.GetCustomerByID(order.CustomerId, token)
	if err != nil {
		log.Errorf("Error getting customer by ID: %v", err)
		return nil
	}
	if customer == nil {
		log.Errorf("Customer not found for ID: %s", order.CustomerId)
		return nil
	}
	log.Infof("Customer Info: %+v", customer)

	priceWithVat := CalculateVat(order.Price)
	log.Infof("Price with VAT: %.2f", priceWithVat)

	// KDV dahil fiyatı güncelle
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
		return nil
	}

	log.Infof("Consumer saved successfully: %+v", consum)
	return nil
}
func CalculateVat(price float64) float64 {
	// KDV oranı
	vatRate := 0.18

	// Fiyat ve KDV oranını çarparak KDV miktarını bulun
	vatAmount := price * vatRate

	// KDV'yi orijinal fiyata ekleyerek toplam fiyatı bulun
	totalPrice := price + vatAmount

	return totalPrice
}
