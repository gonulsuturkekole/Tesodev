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
			return nil, fmt.Errorf("no order found with ID %s", id)
		}
	}
	return order, nil
}

func (r *Repository) Create(ctx context.Context, order interface{}) (*mongo.InsertOneResult, error) {

	res, err := r.collection.InsertOne(ctx, order)
	return res, err
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
