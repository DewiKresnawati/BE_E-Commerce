package repository

import (
	"context"
	"go-loc/config"
	"go-loc/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func FindNearestRoad(lat, lng, maxDistance float64) (*model.Road, error) {
	collection := config.MongoClient.Database("petapedia").Collection("roads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Geospatial Query Operators
	filter := bson.M{
		"geometry": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"$maxDistance": maxDistance,
			},
		},
	}

	// Execute query
	var road model.Road
	err := collection.FindOne(ctx, filter).Decode(&road)
	if err != nil {
		return nil, err
	}
	return &road, nil
}
