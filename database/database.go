package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// LoadEnv loads the first .env found in the working directory or any parent directory.
func LoadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("[Database]: cannot get working directory.")
		return
	}
	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Println("[Database]: failed to load", path)
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Println("[Database]: env files not found!.")
			return
		}
		dir = parent
	}
}

func (conf *Config) connect() *mongo.Database {
	LoadEnv()

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
