package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/OrderService/client"
	"tesodev-korpes/OrderService/internal/types"
	_ "tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/producer"
	"time"
)

type Service struct {
	repo      *Repository
	cusClient *client.CustomerClient
	producer  *producer.Producer
}

func NewService(repo *Repository, cusClient *client.CustomerClient, producer *producer.Producer) *Service {
	return &Service{
		repo:      repo,
		cusClient: cusClient,
		producer:  producer,
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
		Price:            orderReq.Price,
		OrderName:        orderReq.OrderName,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	_, err = s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	err = s.producer.ProduceMessage(order.Id)
	if err != nil {
		log.Printf("Failed to produce orderID to Kafka: %v", err)
	}

	err = s.cusClient.SendFinanceRequest(token)
	if err != nil {
		log.Printf("Failed to send finance request: %v", err)
		// Handle error as needed, perhaps returning a warning or proceeding
		return "", err
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
