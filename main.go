package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors" // Import CORS middleware
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

	// Use logger middleware
	app.Use(logger.New())

	// Use CORS middleware (can customize it in config/cors.go)
	app.Use(cors.New()) // Default CORS settings

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
