package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/CustomerService/internal/types"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Customer, error) {
	var customer *types.Customer
	/*
		return customer, nil*/
	filter := bson.M{"_id": id}

	// Define a variable to hold the result
	//var customer Customer

	// Perform the find operation

	err := r.collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no customer found with ID %s", id)
		}
	}
	return customer, nil
}

// Create method in Repository inserts a customer into MongoDB
func (r *Repository) Create(ctx context.Context, customer *types.Customer) (*mongo.InsertOneResult, error) {
	res, err := r.collection.InsertOne(ctx, customer)
	return res, err
}

func (r *Repository) Update(ctx context.Context, id string, customer *types.Customer) error {
	filter := bson.D{{"_id", id}}
	update := bson.M{"$set": customer}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
