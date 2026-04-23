package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Beer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Brand   string             `bson:"brand" json:"brand"`
	Alcohol float64            `bson:"alcohol" json:"alcohol"`
	Year    int                `bson:"year" json:"year"`
}

const (
	ErrCodeInvalidAlcohol = "INVALID_ALCOHOL_VOLUME"
	ErrCodeInvalidYear    = "INVALID_BREW_YEAR"
	ErrCodeMissingData    = "MISSING_REQUIRED_DATA"
	ErrCodeDuplicateBeer  = "DUPLICATE_BEER"
)
