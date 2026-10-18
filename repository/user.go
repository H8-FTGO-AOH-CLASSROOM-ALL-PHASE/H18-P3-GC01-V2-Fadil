package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID       int64  `bson:"_id" json:"id"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"` // stores the hashed password
}

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{col: db.Collection("users")}
}

// EnsureIndexes creates the unique index on email.
func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// Create - POST /register
func (r *UserRepository) Create(ctx context.Context, u *User) error {
	id, err := nextID(ctx, r.col)
	if err != nil {
		return err
	}
	u.ID = id
	if _, err = r.col.InsertOne(ctx, u); mongo.IsDuplicateKeyError(err) {
		return ErrEmailDuplicate
	} else if err != nil {
		return err
	}
	return nil
}

// FindByEmail - used by POST /login
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.col.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
