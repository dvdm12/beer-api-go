package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Beer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name    string             `bson:"name" json:"name" binding:"required"`
	Brand   string             `bson:"brand" json:"brand" binding:"required"`
	Alcohol float64            `bson:"alcohol" json:"alcohol" binding:"required,gt=0"`
	Year    int                `bson:"year" json:"year" binding:"required,gt=0"`
}
