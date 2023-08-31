package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models/dto"
)

func convertFactToDTO(val *models.Fact) *dto.FactsDTO {
	return &dto.FactsDTO{
		ID:       int(val.ID),
		Question: val.Question,
		Answer:   val.Answer,
	}
}

func convertFactsToDTO(facts []models.Fact) []dto.FactsDTO {

	factsDTO := []dto.FactsDTO{}

	for _, val := range facts {
		dto := convertFactToDTO(&val)
		factsDTO = append(factsDTO, *dto)
	}

	return factsDTO
}

func ListFacts(c *fiber.Ctx) error {
	facts := []models.Fact{}

	database.DB.Db.Find(&facts)
	factsDTO := convertFactsToDTO(facts)
	return c.Status(200).JSON(factsDTO)
}

func GetSingleFact(c *fiber.Ctx) error {
	facts := []models.Fact{}
	id := c.Params("id")

	result := database.DB.Db.Find(&facts, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	factsDTO := convertFactsToDTO(facts)

	return c.Status(200).JSON(factsDTO[0])
}

func CreateFacts(c *fiber.Ctx) error {
	fact := new(models.Fact)

	if err := c.BodyParser(fact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	result := database.DB.Db.Create(&fact)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	dto := convertFactToDTO(fact)

	return c.Status(200).JSON(dto)

}

func DeleteFact(c *fiber.Ctx) error {
	factId := c.Params("id")

	fact := &models.Fact{}
	database.DB.Db.Find(fact, "id = ?", factId)
	if fact.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "No data with that Id exists",
		})
	}

	dto := convertFactToDTO(fact)
	result := database.DB.Db.Delete(fact, "id = ?", factId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(200).JSON(dto)

}
