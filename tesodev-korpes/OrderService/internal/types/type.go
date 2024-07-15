package types

import "time"

type OrderRequestModel struct {
	Id              string  `bson:"_id" json:"id"`
	Price           float64 `bson:"price" json:"price"`
	Stock           int     `bson:"stock" json:"stock"`
	ShippingAddress string  `bson:"shipping_address" json:"shipping_address"`
	PaymentMethod   string  `bson:"payment_method" json:"payment_method"`
}

type OrderResponseModel struct {
	OrderName       string  `bson:"order_name" json:"order_name"`
	Price           float64 `bson:"price" json:"price"`
	ShippingAddress string  `bson:"shipping_address" json:"shipping_address"`
	PaymentMethod   string  `bson:"payment_method" json:"payment_method"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
