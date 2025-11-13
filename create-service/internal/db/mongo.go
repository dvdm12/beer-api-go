package db

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() MongoCollectionInterface {

    uri := os.Getenv("MONGO_URI")
    database := os.Getenv("DATABASE")
    collectionName := os.Getenv("COLLECTION")

    if uri == "" || database == "" || collectionName == "" {
        log.Fatal("Missing required environment variables: MONGO_URI, DATABASE, COLLECTION")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDB ping failed:", err)
    }

    log.Println("MongoDB connected successfully")

    col := client.Database(database).Collection(collectionName)

    return NewMongoCollectionAdapter(col)
}
