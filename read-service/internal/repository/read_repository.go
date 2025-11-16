package repository

import (
	"context"
	"time"

	"readservice/internal/db"
	"readservice/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadRepository struct {
	collection db.MongoCollectionInterface
}

func NewReadRepository(collection db.MongoCollectionInterface) *ReadRepository {
	return &ReadRepository{
		collection: collection,
	}
}

func (r *ReadRepository) GetBeerByID(id string) (*models.Beer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}

	var beer models.Beer
	err = r.collection.FindOne(ctx, filter).Decode(&beer)
	if err != nil {
		return nil, err
	}

	return &beer, nil
}

func (r *ReadRepository) GetAllBeers() ([]models.Beer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var beers []models.Beer
	if err = cursor.All(ctx, &beers); err != nil {
		return nil, err
	}

	return beers, nil
}
