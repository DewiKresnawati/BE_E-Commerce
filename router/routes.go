package router

import (
	"github.com/gofiber/fiber/v2"
	"go-loc/handler"
)

func SetupRoutes(app *fiber.App) {
	roadGroup := app.Group("/roads")
	roadGroup.Post("/nearest", handler.NearestLocation)
}
