package types

import "time"

type Consumer struct {
	Id         string    `bson:"_id" json:"id"`
	CustomerId string    `bson:"customer_id" json:"customer_id"`
	OrderId    string    `bson:"order_id" json:"order_id"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
