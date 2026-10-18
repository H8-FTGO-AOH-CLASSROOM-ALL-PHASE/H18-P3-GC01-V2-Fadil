package repository

import (
	"context"
	"errors"

	"H18-P3-GC01-V2-Fadil/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Product struct {
	ID    int64   `bson:"_id" json:"id"`
	Name  string  `bson:"name" json:"name"`
	Price float64 `bson:"price" json:"price"`
	Qty   int     `bson:"qty" json:"qty"`
}

type ProductRepository struct {
	col *mongo.Collection
}

func NewProductRepository(conf *database.Config) *ProductRepository {
	return &ProductRepository{col: conf.GetNoSQLDatabaseConnection().Collection("products")}
}

// Create - POST /products
func (r *ProductRepository) Create(ctx context.Context, p *Product) error {
	id, err := nextID(ctx, r.col)
	if err != nil {
		return err
	}
	p.ID = id
	_, err = r.col.InsertOne(ctx, p)
	return err
}

// FindAll - GET /products
func (r *ProductRepository) FindAll(ctx context.Context) ([]Product, error) {
	cursor, err := r.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	products := []Product{}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// FindByID - GET /products/{id}
func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*Product, error) {
	var p Product
	err := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Update - PUT /products/{id}
func (r *ProductRepository) Update(ctx context.Context, p *Product) error {
	res, err := r.col.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: p.ID}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "name", Value: p.Name},
			{Key: "price", Value: p.Price},
			{Key: "qty", Value: p.Qty},
		}}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete - DELETE /products/{id}
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.col.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
