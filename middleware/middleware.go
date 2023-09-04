package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/config"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	"github.com/melihcanclk/docker-postgres-go-rest-api/helpers"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthMiddleware(c *fiber.Ctx) error {
	var access_token string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		access_token = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("access_token") != "" {
		access_token = c.Cookies("access_token")
	}

	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	tokenClaims, err := helpers.ValidateToken(access_token, config.AccessTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	ctx := context.TODO()
	userid, err := database.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "Token is invalid or session has expired"})
	}

	var user models.User
	err = database.DB.Db.First(&user, "id = ?", userid).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	fmt.Println("Token UUID: ", tokenClaims.TokenUuid)
	c.Locals("user", helpers.ConvertUserToDTO(&user))
	c.Locals("access_token_uuid", tokenClaims.TokenUuid)

	return c.Next()

}
