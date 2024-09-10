package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
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
func (s *Service) CreateOrderService(ctx context.Context, customerID string, orderReq *types.OrderRequestModel, token string) (string, int, error) {
	customer, err := s.cusClient.GetCustomerByID(customerID, token)
	if err != nil {
		return "", 0, err
	}
	if customer == nil {
		return "", 0, fmt.Errorf("customer not found")
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
		return "", 0, err
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
	writer := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order-topic",
		Balancer: &kafka.LeastBytes{},
	}

	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("OrderID"),
		Value: []byte(orderID),
	})

	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	fmt.Printf("OrderID produced to Kafka: %s\n", orderID)
	return writer.Close()
}

func (s *Service) CalculateDailyTotalOrders(ctx context.Context) error {
	// Tüm müşterilerin siparişlerini say
	customers := []string{"customer1", "customer2", "customer3"} // Burada müşteri kimliklerini dinamik olarak çekebilirsiniz.
	today := time.Now().Format("2006-01-02")
	var dailyOrders []types.CustomerOrders

	for _, customerID := range customers {
		count, err := s.repo.CountOrdersByCustomerID(ctx, customerID)
		if err != nil {
			return fmt.Errorf("failed to count orders for customer %s: %v", customerID, err)
		}

		// Her müşterinin sipariş sayısını kaydet
		dailyOrders = append(dailyOrders, types.CustomerOrders{
			CustomerId: customerID,
			Count:      int(count),
		})
	}

	// Günlük sipariş özetini kaydet
	dailyOrderSummary := types.DailyOrder{
		Date:   today,
		Orders: dailyOrders,
	}

	// MongoDB'ye kaydet
	err := s.repo.SaveDailyOrderSummary(ctx, dailyOrderSummary)
	if err != nil {
		return fmt.Errorf("failed to save daily order summary: %v", err)
	}

	fmt.Println("Günlük sipariş özeti başarıyla kaydedildi.")
	return nil
}
