package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBeer_JSONSerialization(t *testing.T) {
	beer := Beer{
		ID:      primitive.NewObjectID(),
		Name:    "Test Beer",
		Brand:   "Test Brand",
		Alcohol: 4.5,
		Year:    2023,
	}

	jsonData, err := json.Marshal(beer)
	assert.Nil(t, err)
	assert.Contains(t, string(jsonData), "\"name\":\"Test Beer\"")
	assert.Contains(t, string(jsonData), "\"brand\":\"Test Brand\"")
	assert.Contains(t, string(jsonData), "\"alcohol\":4.5")
	assert.Contains(t, string(jsonData), "\"year\":2023")
}

func TestBeer_JSONDeserialization(t *testing.T) {
	jsonStr := `{
		"name": "Colombian Lager",
		"brand": "Aguila",
		"alcohol": 4.0,
		"year": 2022
	}`

	var beer Beer
	err := json.Unmarshal([]byte(jsonStr), &beer)
	assert.Nil(t, err)
	assert.Equal(t, "Colombian Lager", beer.Name)
	assert.Equal(t, "Aguila", beer.Brand)
	assert.Equal(t, 4.0, beer.Alcohol)
	assert.Equal(t, 2022, beer.Year)
}

func TestBeer_BSONSerialization(t *testing.T) {
	beer := Beer{
		ID:      primitive.NewObjectID(),
		Name:    "Stout Pro",
		Brand:   "BogotáBeers",
		Alcohol: 6.2,
		Year:    2021,
	}

	bsonData, err := bson.Marshal(beer)
	assert.Nil(t, err)

	var decoded Beer
	err = bson.Unmarshal(bsonData, &decoded)
	assert.Nil(t, err)

	assert.Equal(t, beer.Name, decoded.Name)
	assert.Equal(t, beer.Brand, decoded.Brand)
	assert.Equal(t, beer.Alcohol, decoded.Alcohol)
	assert.Equal(t, beer.Year, decoded.Year)
}

func TestBeer_DefaultID_Omitted(t *testing.T) {
	beer := Beer{
		Name:    "IPA",
		Brand:   "CraftCo",
		Alcohol: 5.8,
		Year:    2020,
	}

	bsonData, err := bson.Marshal(beer)
	assert.Nil(t, err)

	assert.NotContains(t, string(bsonData), "_id")
}
