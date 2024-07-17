package internal

import (
	"context"
	"github.com/google/uuid"
	"tesodev-korpes/OrderService/internal/types"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	//challenge (everything should be observable somehow in the response or console (print)):
	// 1) do something with using for loop by using customer model and manipulate it (you can add an additional field for it)
	// 2) do something with switch-case
	// 3) do something with goroutines (you should give us an example for both scenarios of not using goroutines and using)
	// 3.1) calculate the elapsed time for both scenarios and show us the gained time
	// 4) add an additional field and use maps
	// 5) add an additional field and use arrays
	// 6) manipulate an existing data to see how pointers and values work
	return order, nil
}

func (s *Service) Create(ctx context.Context, order *types.Order) (string, error) {
	// Generate a new UUID
	orderID := uuid.New().String()
	now := time.Now().Local()

	// Set the customer's ID to the generated UUID
	order.Id = orderID
	order.CreatedAt = now
	// Insert the customer data into MongoDB
	_, err := s.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	// Return the generated ID if the insertion is successful
	return orderID, nil
}

func (s *Service) Update(ctx context.Context, id string, orderUpdateModel types.OrderUpdateModel) error {
	order, err := s.GetByID(ctx, id)
	now := time.Now().Local()
	if err != nil {
		return err
	}

	order.OrderName = orderUpdateModel.OrderName
	order.Price = orderUpdateModel.Price
	order.Stock = orderUpdateModel.Stock
	order.ShippingAddress = orderUpdateModel.ShippingAddress
	order.PaymentMethod = orderUpdateModel.PaymentMethod
	order.UpdatedAt = now
	return s.repo.Update(ctx, id, order)

}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
