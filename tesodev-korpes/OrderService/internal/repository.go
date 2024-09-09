package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/OrderService/internal/types"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Order, error) {
	var order *types.Order

	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no order found with ID %s", id) //nil,nil
		}
	}
	return order, nil //error
}

func (r *Repository) Create(ctx context.Context, order interface{}) (*mongo.InsertOneResult, error) {

	res, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) Update(ctx context.Context, id string, order *types.Order) error {
	filter := bson.D{{"_id", id}}
	update := bson.M{"$set": order}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *Repository) CountOrdersByCustomerID(ctx context.Context, customerID string) (int64, error) {
	filter := bson.M{"customer_id": customerID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count orders for customer %s: %v", customerID, err)
	}
	return count, nil
}

func (r *Repository) SaveDailyOrderSummary(ctx context.Context, dailyOrder types.DailyOrder) error {
	_, err := r.collection.InsertOne(ctx, dailyOrder)
	if err != nil {
		return fmt.Errorf("failed to insert daily order summary: %v", err)
	}
	return nil
}
