package repository

import (
    "context"
    "createservice/internal/db"
    "createservice/internal/models"
    "time"
)

type CreateRepository struct {
    collection db.MongoCollectionInterface
}

func NewCreateRepository(collection db.MongoCollectionInterface) *CreateRepository {
    return &CreateRepository{
        collection: collection,
    }
}

func (r *CreateRepository) CreateBeer(beer models.Beer) (*struct{}, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, beer)
    return nil, err
}
