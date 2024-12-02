package router

import (
	"github.com/gofiber/fiber/v2"
	"go-loc/handler"
)

func SetupRoutes(app *fiber.App) {
	roadGroup := app.Group("/api")
	roadGroup.Post("/getroad", handler.GetRoad)

	// Menambahkan route baru untuk GetRegion
	roadGroup.Post("/getregion", handler.GetRegion)
}
