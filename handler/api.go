package handler

import (
	"context"
	"encoding/json"
	"go-loc/config"
	"go-loc/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NearestLocation(c *fiber.Ctx) error {
	// Mendekodekan body request menjadi RequestBody
	var body RequestBody
	err := json.Unmarshal(c.Body(), &body)
	if err != nil {
		var response model.Response
		response.Status = "Error : Body tidak valid"
		response.Response = err.Error()
		// Mengembalikan response error dengan status 400 Bad Request
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Menghubungkan ke database MongoDB
	Collection := config.MongoClient.Database("petapedia").Collection("roads")
	Ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Geospatial Query Filter
	Filter := bson.M{
		"geometry": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{body.Longitude, body.Latitude},
				},
				"$maxDistance": body.MaxDistance,
			},
		},
	}

	// Melakukan query untuk menemukan road terdekat (menggunakan Find untuk mendapatkan beberapa hasil)
	cursor, err := Collection.Find(Ctx, Filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Tidak ditemukan jalan terdekat"
		respn.Response = err.Error()
		// Mengembalikan response error dengan status 404 Not Found
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(respn)
		}
		// Mengembalikan response error lainnya
		return c.Status(fiber.StatusInternalServerError).JSON(respn)
	}
	defer cursor.Close(Ctx)

	// Membuat slice untuk menampung hasil query
	var roads []model.Roads
	for cursor.Next(Ctx) {
		var road model.Roads
		if err := cursor.Decode(&road); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "Error",
				"message": "Gagal mendekodekan data jalan",
			})
		}
		roads = append(roads, road)
	}

	// Mengembalikan hasil jalan terdekat dalam bentuk JSON
	if len(roads) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "Error",
			"message": "Tidak ada jalan terdekat ditemukan",
		})
	}

	// Membuat struktur GeoJSON yang valid dengan menambahkan "type" dan "features"
	var features []fiber.Map
	for _, road := range roads {
		features = append(features, fiber.Map{
			"type": "Feature", // "type" untuk setiap feature
			"geometry": fiber.Map{
				"type":        road.Geometry.Type,
				"coordinates": road.Geometry.Coordinates,
			},
			"properties": fiber.Map{
				"osm_id":  road.Properties.OSMID,
				"name":    road.Properties.Name,
				"highway": road.Properties.Highway,
			},
		})
	}

	// Menyiapkan response GeoJSON yang valid
	geojsonResponse := fiber.Map{
		"type":     "FeatureCollection", // Menambahkan "type" di tingkat atas
		"features": features,            // Menyertakan array "features" yang berisi data jalan
	}

	// Mengembalikan response dalam format GeoJSON
	return c.JSON(geojsonResponse)
}
