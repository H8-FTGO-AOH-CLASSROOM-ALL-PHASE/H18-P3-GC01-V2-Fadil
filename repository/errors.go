package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound       = errors.New("data not found")
	ErrEmailDuplicate = errors.New("email already registered")
)

// nextID returns the highest _id in the collection plus one.
func nextID(ctx context.Context, col *mongo.Collection) (int64, error) {
	var last struct {
		ID int64 `bson:"_id"`
	}
	err := col.FindOne(ctx, bson.D{}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})).Decode(&last)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	return last.ID + 1, nil
}
