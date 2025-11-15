package repository

import (
    "context"
    "time"

    "updateservice/internal/db"
    "updateservice/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateRepository struct {
    collection db.MongoCollectionInterface
}

func NewUpdateRepository(collection db.MongoCollectionInterface) *UpdateRepository {
    return &UpdateRepository{
        collection: collection,
    }
}

func (r *UpdateRepository) UpdateBeer(id string, beer models.Beer) (*struct{}, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    filter := bson.M{"_id": oid}
    update := bson.M{"$set": beer}

    _, err = r.collection.UpdateOne(ctx, filter, update)
    return nil, err
}
