package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go-loc/config"
	"go-loc/router"
	"log"
	"os"
)

func main() {
	// Initialize MongoDB connection
	config.CreateDBConnection()

	// Initialize Fiber app
	app := fiber.New()

	app.Use(logger.New())

	// Register routes
	router.SetupRoutes(app)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
