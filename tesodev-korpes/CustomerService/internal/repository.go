package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *Repository) GetCustomersByFilter(ctx context.Context, firstName string, ageGreaterThan string, ageLessThan string) ([]types.Customer, error) {
	var customers []types.Customer
	// Create a filter to match the first name
	//opts := options.Find().SetLimit(5)

	//filter := bson.D{
	//{"$or", bson.A{
	//bson.D{{"age", bson.D{{"$gt", age}}}},
	//bson.D{{"first_ame", name}},
	//bson.D{{"last_name", lastName}},
	//}},
	//}
	filter := bson.M{}
	if firstName != "" {
		filter["first_name"] = firstName
	}

	if ageGreaterThan > "" {
		filter["age"] = bson.M{"$gte": ageGreaterThan}
	}
	if ageLessThan > "" {
		if filter["age"] == nil {
			filter["age"] = bson.M{"$lte": ageLessThan}
		} else {
			filter["age"].(bson.M)["$lte"] = ageLessThan
		}
	}
	fmt.Printf("Filter: %v\n", filter) // Log the filter to see what is being sent
	// Perform the query
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "could not get any customers"})
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "error decoding customers"})
	}
	return customers, nil
}
