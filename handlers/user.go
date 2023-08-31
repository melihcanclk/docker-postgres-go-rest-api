package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	"github.com/melihcanclk/docker-postgres-go-rest-api/helpers"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models/dto"
)

func convertUserToDTO(val *models.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:       val.ID,
		Username: val.Username,
		Email:    val.Email,
	}
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error when hashing password", "data": err})
	}

	user.Password = hashed

	result := database.DB.Db.Create(&user)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := convertUserToDTO(user)

	return c.Status(201).JSON(userDTO)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}

	result := database.DB.Db.Find(&user, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := convertUserToDTO(user)

	return c.Status(200).JSON(userDTO)

}

// TODO: Update User
// TODO: Delete User
// TODO: Login
// TODO: Refresh token and bearer token implementation
// https://github.com/adhtanjung/go_rest_api_fiber/blob/main/handler/handler.go