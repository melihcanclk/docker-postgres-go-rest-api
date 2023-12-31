package handlers

import (
	"context"
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

// @Description Create a new user
// @Summary create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserLoginBodyDTO true "User"
// @Success 201 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/register [post]
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
	userDTO := helpers.ConvertUserToDTO(user)

	return c.Status(fiber.StatusCreated).JSON(userDTO)
}

// @Description Login a user, returns access token, refresh token as cookies, user data and access token as json
// @Summary login a user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserLoginBodyDTO true "User"
// @Success 202 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/login [post]
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
	refreshTokenDuration, err := time.ParseDuration(config.RefreshTokenExpiredInMinutes)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	refreshTokenDetails, err := helpers.GenerateJWTToken(user.ID.String(), &refreshTokenDuration, config.RefreshTokenPrivateKey)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	ctx := context.TODO()

	// save refresh token to redis
	err = database.RedisClient.Set(ctx, refreshTokenDetails.TokenUuid, user.ID.String(), refreshTokenDuration).Err()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// save access token to redis
	err = database.RedisClient.Set(ctx, accessTokenDetails.TokenUuid, user.ID.String(), accessTokenDuration).Err()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	accessTokenMaxAge := int(config.AccessTokenMaxAge)
	expiredDay := time.Now().Add(time.Minute * time.Duration(accessTokenMaxAge))

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   accessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
		Expires:  expiredDay,
	})

	refreshTokenMaxAge := int(config.RefreshTokenMaxAge)
	expiredDay = time.Now().Add(time.Minute * time.Duration(refreshTokenMaxAge))

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		MaxAge:   refreshTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
		Expires:  expiredDay,
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

	userDTO := helpers.ConvertUserToDTO(user)

	c.Locals("user", userDTO)
	c.Locals("access_token_uuid", accessTokenDetails.TokenUuid)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":       "success",
		"message":      "Login Success",
		"user":         userDTO,
		"access_token": accessTokenDetails.Token,
	})
}

// @Description Logout a user, deletes refresh token from redis and access token from cookie
// @Summary logout a user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /auth/v1/logout [get]
func LogoutUser(c *fiber.Ctx) error {

	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "Token is invalid or session has expired"})
	}

	ctx := context.TODO()

	tokenClaims, err := helpers.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	access_token_uuid := c.Locals("access_token_uuid").(string)
	_, err = database.RedisClient.Del(ctx, tokenClaims.TokenUuid, access_token_uuid).Result()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// delete access token from cookie
	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: "",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: "",
	})

	c.Cookie(&fiber.Cookie{
		Name:  "logged_in",
		Value: "false",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Successfully logged out"})

}

// @Description Get a user with given id
// @Summary get a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}

	result := database.DB.Db.Find(&user, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := helpers.ConvertUserToDTO(user)

	return c.Status(fiber.StatusOK).JSON(userDTO)
}

// @Description Update a user with given id
// @Summary update a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UserUpdateBodyDTO true "User"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/users/{id} [put]
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
	userDTO := helpers.ConvertUserToDTO(user)
	return c.Status(fiber.StatusOK).JSON(userDTO)
}

// @Description Delete a user with given id
// @Summary delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}

	result := database.DB.Db.Find(&user, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	userDTO := helpers.ConvertUserToDTO(user)
	database.DB.Db.Delete(&user, "id = ?", id)

	return c.Status(fiber.StatusOK).JSON(userDTO)
}

// @Description Get current user
// @Summary get current user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/users/me [get]
func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.UserDTO)
	fmt.Println(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

// @Description Refresh access token
// @Summary refresh access token
// @Tags Users
// @Accept json
// @Produce json
// @Success 202 {object} dto.UserDTO
// @Failure 400 {object} string
// @Router /auth/v1/refresh [get]
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

	// get user from redis
	ctx := context.TODO()
	userID, err := database.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// get user from database
	user := &models.User{}
	err = database.DB.Db.Find(&user, "id = ?", userID).Error

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

	// save access token to redis
	err = database.RedisClient.Set(ctx, accessTokenDetails.TokenUuid, user.ID, accessTokenDuration).Err()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	accessTokenMaxAge := int(config.AccessTokenMaxAge)
	expiredDay := time.Now().Add(time.Minute * time.Duration(accessTokenMaxAge))

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   accessTokenMaxAge,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
		Expires:  expiredDay,
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

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":       "success",
		"message":      "Access token refreshed",
		"access_token": accessTokenDetails.Token,
	})
}
