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

// @Description List all facts
// @Summary get all facts
// @Tags Facts
// @Accept json
// @Produce json
// @Success 200 {object} []dto.FactsDTO
// @Failure 404 {object} string
// @Router /api/v1/facts [get]
func ListFacts(c *fiber.Ctx) error {
	facts := []models.Fact{}

	database.DB.Db.Find(&facts)
	factsDTO := convertFactsToDTO(facts)
	return c.Status(200).JSON(factsDTO)
}

// @Description Get a single fact
// @Summary get a single fact
// @Tags Facts
// @Accept json
// @Produce json
// @Param id path int true "Fact ID"
// @Success 200 {object} dto.FactsDTO
// @Failure 404 {object} string
// @Router /api/v1/facts/{id} [get]
func GetSingleFact(c *fiber.Ctx) error {
	fact := &models.Fact{}
	id := c.Params("id")

	// get first fact matching the id
	result := database.DB.Db.Find(&fact, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	factsDTO := convertFactToDTO(fact)

	return c.Status(200).JSON(factsDTO)
}

// @Description Create a fact
// @Summary create a fact
// @Tags Facts
// @Accept json
// @Produce json
// @Param question body string true "Question"
// @Param answer body string true "Answer"
// @Success 200 {object} dto.FactsDTO
// @Failure 404 {object} string
// @Router /api/v1/facts [post]
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

// @Description Delete a fact
// @Summary delete a fact
// @Tags Facts
// @Accept json
// @Produce json
// @Param id path int true "Fact ID"
// @Success 200 {object} dto.FactsDTO
// @Failure 404 {object} string
// @Router /api/v1/facts/{id} [delete]
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
