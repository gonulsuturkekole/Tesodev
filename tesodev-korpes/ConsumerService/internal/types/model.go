package types

import (
	"time"
)

type Consumer struct {
	Id        string                `bson:"_id" json:"id"`
	Customer  CustomerResponseModel `bson:"customer" json:"customer"`
	Order     OrderResponseModel    `bson:"order" json:"order"`
	CreatedAt time.Time             `bson:"created_at" json:"created_at"`
}
