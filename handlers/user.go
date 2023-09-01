package handlers

import (
	"strings"

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

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	body, user := &dto.UserUpdateBodyEntity{}, &models.User{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "Error when parsing body"})
	}
	result := database.DB.Db.Find(&user, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	if body.Username != "" {
		body.Username = strings.ToLower(body.Username)
		err := helpers.IsIncludesNonAscii(&body.Username)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		user.Username = body.Username

	}
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		hashed, err := helpers.HashPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error when hashing password", "data": err})
		}
		user.Password = hashed
	}

	result = database.DB.Db.Save(&user)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := convertUserToDTO(user)
	return c.Status(200).JSON(userDTO)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}

	result := database.DB.Db.Delete(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	userDTO := convertUserToDTO(user)
	return c.Status(200).JSON(userDTO)

}

// TODO: Delete User
// TODO: Login
// TODO: Refresh token and bearer token implementation
// https://github.com/adhtanjung/go_rest_api_fiber/blob/main/handler/handler.go
