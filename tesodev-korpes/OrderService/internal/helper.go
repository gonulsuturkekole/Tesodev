package internal

import (
	"tesodev-korpes/OrderService/internal/types"
)

func ToOrderResponse(order *types.Order) *types.OrderResponseModel {
	return &types.OrderResponseModel{
		CustomerId:     order.CustomerId,
		OrderName:      order.OrderName,
		ShipmentStatus: order.ShipmentStatus,
		PaymentMethod:  order.PaymentMethod,
		OrderTotal:     order.OrderTotal,
		Price:          order.Price,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

}
