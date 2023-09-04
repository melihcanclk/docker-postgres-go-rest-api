package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
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
	v1.Get("/logout", middleware.AuthMiddleware, handlers.LogoutUser)

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware)
	users.Get("/me", handlers.GetMe)
	users.Get("/:id", handlers.GetUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

}

func setupSwaggerRoutes(app *fiber.App) {

	swaggerRoute := app.Group("/swagger")
	// swaggerRoute.Get("/*", swagger.HandlerDefault) // default

	swaggerRoute.Get("/*", swagger.New(swagger.Config{ // custom
		URL:          "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		OAuth2RedirectUrl: "http://localhost:3000/swagger/oauth2-redirect.html",
	}))
}
