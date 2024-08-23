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
	repo          *FinanceRepository
	kafkaConsumer *consumer.Consumer
	conClient     *clientCon.ConsumerClient
}

func NewService(repo *FinanceRepository, conClient *clientCon.ConsumerClient, kafkaConsumer *consumer.Consumer, brokers []string, topic string) *Service {

	kafkaConsumer.Topic = topic
	kafkaConsumer.CreateConnection(brokers)

	return &Service{
		repo:          repo,
		conClient:     conClient,
		kafkaConsumer: kafkaConsumer,
	}
}

func (s *Service) ProcessMessage(ctx context.Context, msg string, key string) error {

	err := s.aggregateCustomerOrder(ctx, msg, key)
	if err != nil {
		fmt.Printf("Error sending order request: %v\n", err)
		return err
	}
	return nil
}

func (s *Service) aggregateCustomerOrder(ctx context.Context, msg string, key string) error {

	order, err := s.conClient.GetOrderByID(msg, key)
	if err != nil {
		log.Errorf("Error getting order by ID: %v", err)
		return nil
	}
	if order == nil {
		log.Errorf("Order not found for ID: %s", msg)
		return nil
	}
	log.Infof("Order Info: %+v", order)

	customer, err := s.conClient.GetCustomerByID(order.CustomerId, key)
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

	order.Price = priceWithVat

	consum := &types.CustomerOrder{
		Id:       uuid.New().String(),
		Customer: *customer,
		Order:    *order,
	}

	_, err = s.repo.Create(ctx, consum)
	if err != nil {
		log.Errorf("Error saving consumer to repository: %v", err)
		return nil
	}

	log.Infof("CustomerOrder saved successfully: %+v", consum)
	return nil
}
func CalculateVat(price float64) float64 {

	vatRate := 0.18
	vatAmount := price * vatRate
	totalPrice := price + vatAmount

	return totalPrice
}
