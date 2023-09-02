package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/handlers"
	"github.com/melihcanclk/docker-postgres-go-rest-api/middleware"
)

func setupFactsRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(middleware.AuthMiddleware)

	v1.Get("/facts", handlers.ListFacts)
	v1.Get("/facts/:id", handlers.GetSingleFact)
	v1.Post("/facts", handlers.CreateFacts)
	v1.Delete("/facts/:id", handlers.DeleteFact)
}

func setupUserRoutes(app *fiber.App) {
	api := app.Group("/auth")
	v1 := api.Group("/v1")

	v1.Post("/register", handlers.CreateUser)
	v1.Post("/login", handlers.LoginUser)
	v1.Get("/refresh", handlers.RefreshAccessToken)
	//TODO: logout

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware)
	users.Get("/me", handlers.GetMe)
	users.Get("/:id", handlers.GetUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

}
