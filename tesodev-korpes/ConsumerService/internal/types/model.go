package types

import (
	"time"
)

type CustomerOrder struct {
	Id string `bson:"_id" json:"id"`
	//Customer  CustomerResponseModel `bson:"customer" json:"customer"`
	FirstName      string    `bson:"first_name" json:"first_name"`
	LastName       string    `bson:"last_name" json:"last_name"`
	Username       string    `bson:"username" json:"username"`
	CustomerId     string    `bson:"customer_id" json:"customer_id"`
	OrderName      string    `bson:"order_name" json:"order_name"`
	ShipmentStatus string    `bson:"shipment_status" json:"shipment_status"`
	PaymentMethod  string    `bson:"payment_method" json:"payment_method"`
	OrderTotal     int       `bson:"order_total" json:"order_total"`
	PriceCent      int       `bson:"priceCent" json:"priceCent"`
	OrderCreatedAt time.Time `bson:"order_created_at" json:"order-created-at"`
	OrderUpdatedAt time.Time `bson:"order-updated-at" json:"order-updated-at"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}
