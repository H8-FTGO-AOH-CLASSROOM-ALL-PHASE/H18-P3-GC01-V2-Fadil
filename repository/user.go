package repository

import (
	"context"
	"errors"

	"H18-P3-GC01-V2-Fadil/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID        int64  `bson:"_id" json:"id"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"-"`   // stores the hashed password
	JWTSecret string `bson:"jwt_secret" json:"-"` // unique signing secret for this user
}

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(conf *database.Config) *UserRepository {
	return &UserRepository{col: conf.GetNoSQLDatabaseConnection().Collection("users")}
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

// FindByID - used by the JWT middleware to load the user's secret
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
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
