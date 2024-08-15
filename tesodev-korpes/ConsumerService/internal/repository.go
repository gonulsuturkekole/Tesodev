package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}
func (r *Repository) ConsumeTopic(ctx context.Context) error {
	// Placeholder for future implementation
	return nil
}

func (r *Repository) Create(ctx context.Context, consum interface{}) (*mongo.InsertOneResult, error) {

	res, err := r.collection.InsertOne(ctx, consum)
	if err != nil {
		return nil, err
	}
	return res, nil
}

/*func (r *Repository) FindByID(ctx context.Context, id string) (*types.Order, error) {
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
*/
