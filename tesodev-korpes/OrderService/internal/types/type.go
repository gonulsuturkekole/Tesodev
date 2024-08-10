package types

import "time"

type OrderRequestModel struct {
	CustomerId    string `bson:"customer_id" json:"customer_id"`
	OrderTotal    int    `bson:"order_total" json:"order_total"`
	PaymentMethod string `bson:"payment_method" json:"payment_method"`
}

type OrderResponseModel struct {
	CustomerId     string `bson:"customer_id" json:"customer_id"`
	OrderName      string `bson:"order_name" json:"order_name"`
	ShipmentStatus string `bson:"shipment_status" json:"shipment_status"`
	PaymentMethod  string `bson:"payment_method" json:"payment_method"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
type OrderUpdateModel struct {
	OrderName      string    `bson:"order_name" json:"order_name"`
	ShipmentStatus string    `bson:"shipment_status" json:"shipment_status"`
	PaymentMethod  string    `bson:"payment_method" json:"payment_method"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

type CustomerResponse struct {
	UserName string `bson:"user_name" json:"user_name"`
	Name     string `bson:"name" json:"name"`
	LastName string `bson:"last_name" json:"last_name"`
}
