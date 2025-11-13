package repository

import (
    "context"
    "os"
    "time"

    "createservice/internal/models"
    "go.mongodb.org/mongo-driver/mongo"
)

type CreateRepository struct {
    collection *mongo.Collection
}

func NewCreateRepository(client *mongo.Client) *CreateRepository {
    database := os.Getenv("DATABASE")
    collection := os.Getenv("COLLECTION")

    return &CreateRepository{
        collection: client.Database(database).Collection(collection),
    }
}

func (r *CreateRepository) CreateBeer(beer models.Beer) (*mongo.InsertOneResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return r.collection.InsertOne(ctx, beer)
}
