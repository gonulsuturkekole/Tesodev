package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/ConsumerService/internal/types"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}
func (r *Repository) Create(ctx context.Context, customer *types.Consumer) (*mongo.InsertOneResult, error) {
	res, err := r.collection.InsertOne(ctx, customer)
	return res, err
}
