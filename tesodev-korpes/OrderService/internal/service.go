package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/OrderService/client"
	"tesodev-korpes/OrderService/internal/types"
	_ "tesodev-korpes/pkg"
	"time"
)

type Service struct {
	repo      *Repository
	cusClient *client.CustomerClient
}

func NewService(repo *Repository, cusClient *client.CustomerClient) *Service {
	return &Service{
		repo:      repo,
		cusClient: cusClient,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// if order nil (404 error)

	return order, nil
}

func (s *Service) CreateOrderService(ctx context.Context, customerID string, orderReq *types.OrderRequestModel, token string) (string, error) {

	customer, err := s.cusClient.GetCustomerByID(customerID, token)
	if err != nil {
		return "", err
	}
	if customer == nil {
		return "", fmt.Errorf("customer not found")
	}

	now := time.Now().Local()

	order := &types.Order{
		Id:               uuid.New().String(),
		CustomerId:       customerID,
		CustomerResponse: *customer,
		OrderTotal:       orderReq.OrderTotal,
		PaymentMethod:    orderReq.PaymentMethod,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	_, err = s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	err = s.produceToKafka(order.Id)
	if err != nil {
		log.Printf("Failed to produce orderID to Kafka: %v", err)
	}
	return order.Id, nil
}
func (s *Service) Update(ctx context.Context, id string, orderUpdateModel types.OrderUpdateModel) error {
	order, err := s.GetByID(ctx, id)
	now := time.Now().Local()
	if err != nil {
		return err
	}

	order.OrderName = orderUpdateModel.OrderName
	order.ShipmentStatus = orderUpdateModel.ShipmentStatus
	order.PaymentMethod = orderUpdateModel.PaymentMethod
	order.UpdatedAt = now
	return s.repo.Update(ctx, id, order)

}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) produceToKafka(orderID string) error {

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: "order-topic", NumPartitions: 1, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order-topic",
		Balancer: &kafka.LeastBytes{},
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("OrderID"),
		Value: []byte(orderID),
	})

	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	fmt.Printf("OrderID produced to Kafka: %s\n", orderID)
	return writer.Close()
}
