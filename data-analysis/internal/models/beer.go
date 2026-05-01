package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Beer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Brand   string             `bson:"brand" json:"brand"`
	Alcohol float64            `bson:"alcohol" json:"alcohol"`
	Year    int                `bson:"year" json:"year"`
}

type GeneralStats struct {
	Total       int64   `json:"total"`
	TotalBrands int64   `json:"total_brands"`
	AvgAlcohol  float64 `json:"avg_alcohol"`
	MinAlcohol  float64 `json:"min_alcohol"`
	MaxAlcohol  float64 `json:"max_alcohol"`
	OldestYear  int     `json:"oldest_year"`
	NewestYear  int     `json:"newest_year"`
}

type BrandStats struct {
	Brand        string  `json:"brand"`
	Count        int64   `json:"count"`
	AvgAlcohol   float64 `json:"avg_alcohol"`
	MinAlcohol   float64 `json:"min_alcohol"`
	MaxAlcohol   float64 `json:"max_alcohol"`
	TotalAlcohol float64 `json:"-"`
}
