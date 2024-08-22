package types

import (
	"time"
)

type OrderResponseModel struct {
	CustomerId     string    `bson:"customer_id" json:"customer_id"`
	OrderName      string    `bson:"order_name" json:"order_name"`
	ShipmentStatus string    `bson:"shipment_status" json:"shipment_status"`
	PaymentMethod  string    `bson:"payment_method" json:"payment_method"`
	OrderTotal     int       `bson:"order_total" json:"order_total"`
	Price          float64   `bson:"price" json:"price"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

type CustomerResponseModel struct {
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Username  string `bson:"username" json:"username"`
	Password  string `bson:"password" json:"password"`
}
