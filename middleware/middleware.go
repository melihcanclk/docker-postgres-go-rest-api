package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/config"

	jwtware "github.com/gofiber/jwt/v2"
)

func AuthMiddleware() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Secret),
	})
}
