package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/ConsumerService/internal/types"
)

type FinanceRepository struct {
	collection *mongo.Collection
}

func NewFinanceRepository(col *mongo.Collection) *FinanceRepository {
	return &FinanceRepository{
		collection: col,
	}
}
func (r *FinanceRepository) Create(ctx context.Context, customerOrder *types.CustomerOrder) (*mongo.InsertOneResult, error) {
	res, err := r.collection.InsertOne(ctx, customerOrder)
	return res, err
}
