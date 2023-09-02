package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/config"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	"github.com/melihcanclk/docker-postgres-go-rest-api/helpers"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models/dto"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	user.Username = strings.ToLower(user.Username)
	err := helpers.IsIncludesNonAscii(&user.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
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
	userDTO := ConvertUserToDTO(user)

	return c.Status(fiber.StatusCreated).JSON(userDTO)
}

func LoginUser(c *fiber.Ctx) error {
	body := &dto.UserLoginBodyDTO{}
	user := &models.User{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	var query string
	var value string
	if helpers.IsEmailValid(body.Email) {
		query = "email"
		value = body.Email
	} else {
		query = "username"
		value = body.Username
	}
	result := database.DB.Db.Where(query+" = ?", value).First(&user)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that " + query + " exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	accessTokenDuration, err := time.ParseDuration(config.AccessTokenExpiredInMinutes)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	accessTokenDetails, err := helpers.GenerateJWTToken(user.ID.String(), &accessTokenDuration, config.AccessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	fmt.Println("ref2:", config.RefreshTokenExpiredInMinutes)
	refreshTokenDuration, err := time.ParseDuration(config.RefreshTokenExpiredInMinutes)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	refreshTokenDetails, err := helpers.GenerateJWTToken(user.ID.String(), &refreshTokenDuration, config.RefreshTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	//TODO: these will be used for redis cache after implementation
	// ctx := context.TODO()
	// now := time.Now()

	accessTokenMaxAge := int(config.AccessTokenMaxAge) * 60
	fmt.Println("accmax:", accessTokenMaxAge)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   accessTokenMaxAge,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	refreshTokenMaxAge := int(config.RefreshTokenMaxAge) * 60
	fmt.Println("refmax:", refreshTokenMaxAge)

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		MaxAge:   int(config.RefreshTokenMaxAge) * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	userDTO := ConvertUserToDTO(user)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":       "success",
		"message":      "Login Success",
		"user":         userDTO,
		"access_token": accessTokenDetails.Token,
	})
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
	userDTO := ConvertUserToDTO(user)

	return c.Status(200).JSON(userDTO)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	body, user := &dto.UserUpdateBodyDTO{}, &models.User{}

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
	userDTO := ConvertUserToDTO(user)
	return c.Status(fiber.StatusOK).JSON(userDTO)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}

	result := database.DB.Db.Find(&user, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := ConvertUserToDTO(user)
	database.DB.Db.Delete(&user, "id = ?", id)

	return c.Status(fiber.StatusOK).JSON(userDTO)
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.UserDTO)
	fmt.Println(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func RefreshAccessToken(c *fiber.Ctx) error {
	message := "could not refresh access token"

	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
	}
	tokenClaims, err := helpers.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	err = database.DB.Db.First(&user, "id = ?", tokenClaims.UserID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
		} else {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})

		}
	}

	accessTokenDuration, err := time.ParseDuration(config.AccessTokenExpiredInMinutes)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	accessTokenDetails, err := helpers.GenerateJWTToken(user.ID.String(), &accessTokenDuration, config.AccessTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	accessTokenMaxAge := int(config.AccessTokenMaxAge) * 60

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   accessTokenMaxAge,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   accessTokenMaxAge,
		Secure:   false,
		HTTPOnly: false,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "access_token": accessTokenDetails.Token})
}

func ConvertUserToDTO(val *models.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:       val.ID,
		Username: val.Username,
		Email:    val.Email,
	}
}
