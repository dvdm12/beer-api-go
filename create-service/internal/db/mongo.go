// Package db provides MongoDB connection utilities.
package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect initializes a MongoDB connection and returns a collection adapter.
func Connect() MongoCollectionInterface {

	// Read required environment variables.
	uri := os.Getenv("MONGO_URI")
	database := os.Getenv("DATABASE")
	collectionName := os.Getenv("COLLECTION")

	// Validate required configuration.
	if uri == "" || database == "" || collectionName == "" {
		log.Fatal("Missing required environment variables: MONGO_URI, DATABASE, COLLECTION")
	}

	// Create a context with timeout for connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB client.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Verify connection with a ping.
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	// Log successful connection.
	log.Println("MongoDB connected successfully")

	// Get the target collection.
	col := client.Database(database).Collection(collectionName)

	// Return collection wrapped in adapter.
	return NewMongoCollectionAdapter(col)
}