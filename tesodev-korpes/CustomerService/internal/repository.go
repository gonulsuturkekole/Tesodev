package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
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
func (r *Repository) GetCustomersWithSecondLetterA(ctx context.Context) ([]types.Customer, error) {
	filter := bson.M{"firstName": bson.M{"$regex": "^.{1}a"}}
	opts := options.Find().SetLimit(5)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []types.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
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

func (r *Repository) Get(ctx context.Context) ([]types.Customer, error) {
	var customerModels []types.Customer
	//
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {

		return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "could not get any customers"})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &customerModels); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "error decoding customers"})
	}

	return customerModels, nil
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

//
