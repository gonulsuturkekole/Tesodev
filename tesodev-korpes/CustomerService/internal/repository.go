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

type Pager struct {
	Page      int
	Limit     int
	Offset    int
	AllRecord int
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

/*
		func (r *Repository) GetCustomersByFilter(ctx context.Context, firstName string, ageGreaterThan string, ageLessThan string, page int, limit int) ([]types.Customer,int64, error) {
			var customers []types.Customer
			// Create a filter to match the first name
			pager := NewPager(page, limit)
			filter := bson.M{}
			if firstName != "" {
				filter["first_name"] = firstName
			}
			if ageGreaterThan > "" {
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
			totalCount, err := r.collection.CountDocuments(ctx, filter)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "error counting customers"})
			}
			pager.AllRecord = int(totalCount)

			// Calculate the offset based on the page and limit
			offset := (page - 1) * limit
			opts := options.Find()
			opts.SetSkip(int64(offset))
			opts.SetLimit(int64(limit))

			cursor, err := r.collection.Find(ctx, filter, opts)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "could not get any customers"})
			}
			defer cursor.Close(ctx)
			if err := cursor.All(ctx, &customers); err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "error decoding customers"})
			}
			/*response := map[string]interface{}{
				"customers":  customers,
				"totalCount": totalCount,
			}
			return response, nil//

			return customers, nil
		}
*/
func NewPager(page int, limit int) Pager {
	pager := Pager{}

	if page < 1 {
		pager.Page = 1
	} else {
		pager.Page = page
	}

	if limit <= 0 {

		return Pager{Page: pager.Page, Limit: 0, AllRecord: 0}
	}

	pager.Limit = limit
	pager.Offset = (pager.Page - 1) * pager.Limit
	return pager
}

/*
	func (p *Pager) GetOffset() int {
		return (p.Page - 1) * p.Limit
	}
*/
func (r *Repository) GetCustomersByFilter(ctx context.Context, firstName string, ageGreaterThan string, ageLessThan string, page int, limit int) ([]types.Customer, int64, error) {
	var customers []types.Customer
	pager := NewPager(page, limit)

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
	pager.AllRecord = int(totalCount)

	offset := (page - 1) * limit
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
