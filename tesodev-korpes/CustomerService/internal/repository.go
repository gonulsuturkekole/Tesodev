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
	opts := options.Find().SetLimit(5)
	//filter := bson.D{
	//{"$or", bson.A{
	//bson.D{{"age", bson.D{{"$gt", age}}}},
	//bson.D{{"first_ame", name}},
	//bson.D{{"last_name", lastName}},
	//}},
	//}
	//filters

	/*filter := bson.M{
		"$or": []bson.M{
			{"first_name": firstName},
			{"$and":[{"age": bson.M{"$gt": ageGreaterThan}},{"age": bson.M{"$lt": ageLessThan}}
			]},
		},
	} */
	// Initialize an empty filter map
	filter := bson.M{}
	// Check if firstName is not empty
	if firstName != "" {
		// Add "first_name" to the filter with the value of firstName
		filter["first_name"] = firstName
	}
	// Check if ageGreaterThan is not empty
	if ageGreaterThan > "" {
		// Add "age" to the filter with a condition that it should be greater than or equal to ageGreaterThan
		filter["age"] = bson.M{"$gte": ageGreaterThan}
	}
	// Check if ageLessThan is not empty
	if ageLessThan > "" {
		// Check if "age" is not already in the filter
		if filter["age"] == nil {
			// Add "age" to the filter with a condition that it should be less than or equal to ageLessThan
			filter["age"] = bson.M{"$lte": ageLessThan}
		} else {
			// If "age" is already in the filter, add the less than or equal condition to the existing "age" filter
			filter["age"].(bson.M)["$lte"] = ageLessThan
		}
	}
	fmt.Printf("Filter: %v\n", filter) // Log the filter to see what is being sent
	// Perform the query
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "could not get any customers"})
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "error decoding customers"})
	}
	return customers, nil
}
