package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	setupFactsRoutes(app)
	setupUserRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(fmt.Sprintf(":%d", 3000))
}
