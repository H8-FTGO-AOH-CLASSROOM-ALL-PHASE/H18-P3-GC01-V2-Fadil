package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	DBMongoURI  string
	DBMongoName string

	once sync.Once
	db   *mongo.Database
}

// GetNoSQLDatabaseConnection connects once and reuses the same database handle on later calls.
func (conf *Config) GetNoSQLDatabaseConnection() *mongo.Database {
	conf.once.Do(func() { conf.db = conf.connect() })
	return conf.db
}

func (conf *Config) connect() *mongo.Database {
	//	Load environment variable from .env
	if err := godotenv.Load(); err != nil {
		log.Println("[Database]: env files not found!.")
	}

	//Get MongoDB URI
	mongoDBURI := os.Getenv("MONGODB_URI")
	if mongoDBURI == "" {
		log.Fatal("[MongoDB]: URI not set!.")
	}
	opts := options.Client().ApplyURI(mongoDBURI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client.Database(conf.DBMongoName)
}
