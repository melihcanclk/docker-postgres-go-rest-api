package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	_ "github.com/melihcanclk/docker-postgres-go-rest-api/docs"
)

func initialize(app *fiber.App) {

	database.ConnectDB()
	database.ConnectToRedis()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	setupFactsRoutes(app)
	setupUserRoutes(app)
	setupSwaggerRoutes(app)

	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "No Such Query",
		})
	})

}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()
	initialize(app)

	app.Listen(fmt.Sprintf(":%d", 3000))
}
