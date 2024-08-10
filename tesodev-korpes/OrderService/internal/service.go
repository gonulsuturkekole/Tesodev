package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

	order := &types.Order{
		Id:            uuid.New().String(),
		CustomerId:    customerID,
		OrderTotal:    orderReq.OrderTotal,
		PaymentMethod: orderReq.PaymentMethod,
	}

	orderID := uuid.New().String()
	now := time.Now().Local()
	order.Id = orderID
	order.CreatedAt = now
	order.UpdatedAt = now

	_, err = s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}

	return orderID, nil
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
