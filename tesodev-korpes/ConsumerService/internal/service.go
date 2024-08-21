package internal

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"tesodev-korpes/ConsumerService/clientCon"
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

func (s *Service) ProcessMessage(msg string) error {
	// Mesajı işleyin ve Order Service'e istek gönderin
	err := s.sendRequest(msg)
	if err != nil {
		fmt.Printf("Error sending order request: %v\n", err)
		return err
	}
	return nil
}

func (s *Service) sendRequest(msg string) error {

	req := http.Header{}
	token := req.Get("Authentication")
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

	/*priceWithVat := s.CalculateVat(order.Price)
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

	log.Infof("Consumer saved successfully: %+v", consum)*/
	return nil
}
