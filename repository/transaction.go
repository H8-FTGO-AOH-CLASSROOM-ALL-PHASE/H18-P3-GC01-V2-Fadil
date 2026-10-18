package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Transaction struct {
	ID   int64  `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type TransactionRepository struct {
	col *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{col: db.Collection("transactions")}
}

// Create - POST /transactions
func (r *TransactionRepository) Create(ctx context.Context, t *Transaction) error {
	id, err := nextID(ctx, r.col)
	if err != nil {
		return err
	}
	t.ID = id
	_, err = r.col.InsertOne(ctx, t)
	return err
}

// FindAll - GET /transactions
func (r *TransactionRepository) FindAll(ctx context.Context) ([]Transaction, error) {
	cursor, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	transactions := []Transaction{}
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

// FindByID - GET /transactions/{id}
func (r *TransactionRepository) FindByID(ctx context.Context, id int64) (*Transaction, error) {
	var t Transaction
	err := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&t)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Update - PUT /transactions/{id}
func (r *TransactionRepository) Update(ctx context.Context, t *Transaction) error {
	res, err := r.col.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: t.ID}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: t.Name}}}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete - DELETE /transactions/{id}
func (r *TransactionRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.col.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
