package internal

import (
	"tesodev-korpes/OrderService/internal/types"
)

func ToOrderResponse(order *types.Order) *types.OrderResponseModel {
	return &types.OrderResponseModel{
		OrderName:       order.OrderName,
		Price:           order.Price,
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   order.PaymentMethod,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}

}
