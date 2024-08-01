package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
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

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

//	func (r *Repository) GetUser(ctx context.Context, username string) (*types.CustomerResponseModel, error) {
//		var user types.CustomerResponseModel
//		err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
//		if err != nil {
//			if err == mongo.ErrNoDocuments {
//				return nil, nil
//			}
//			return nil, err
//		}
//		return &user, nil
//	}
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

	filter := bson.M{"_id": id}

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

func NewPager(page int, limit int) Pagination {
	pagination := Pagination{}

	if page < 1 {
		pagination.Page = 1
	} else {
		pagination.Page = page
	}

	if limit <= 0 {
		pagination.Limit = 0
	} else {
		pagination.Limit = limit
	}

	return pagination
}

func (r *Repository) GetCustomersByFilter(ctx context.Context, firstName string, ageGreaterThan string, ageLessThan string, page int, limit int) ([]types.Customer, int64, error) {
	var customers []types.Customer
	pagination := NewPager(page, limit)

	filter := bson.M{}
	if firstName != "" {
		filter["first_name"] = firstName
	}

	// Convert age parameters to integers
	if ageGreaterThan != "" {
		ageGte, err := strconv.Atoi(ageGreaterThan)
		if err != nil {
			return nil, 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid age_greater_than parameter"})
		}
		filter["age"] = bson.M{"$gte": ageGte}
	}

	if ageLessThan != "" {
		ageLte, err := strconv.Atoi(ageLessThan)
		if err != nil {
			return nil, 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid age_less_than parameter"})
		}
		if filter["age"] == nil {
			filter["age"] = bson.M{"$lte": ageLte}
		} else {
			filter["age"].(bson.M)["$lte"] = ageLte
		}
	}

	fmt.Printf("Filter: %v\n", filter)

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error counting customers"})
	}

	offset := (pagination.Page - 1) * pagination.Limit
	pagination.Offset = offset

	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Could not get any customers"})
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error decoding customers"})
	}

	return customers, totalCount, nil

}
