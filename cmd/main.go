package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
)

func initialize(app *fiber.App) {

	database.ConnectDB()
	database.ConnectToRedis()

	app.Use(cors.New())
	app.Use(logger.New())
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "No Such Query",
		})
	})
	setupFactsRoutes(app)
	setupUserRoutes(app)
}

func main() {
	app := fiber.New()
	initialize(app)

	app.Listen(fmt.Sprintf(":%d", 3000))
}
