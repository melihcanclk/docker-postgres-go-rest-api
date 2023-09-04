package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/business"
	"github.com/melihcanclk/docker-postgres-go-rest-api/database"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models/dto"
)

func convertAnswerToDTO(val *[]models.Answer) []dto.AnswersDTO {
	// traverse through the answers
	answersDTO := []dto.AnswersDTO{}
	for _, val := range *val {
		answersDTO = append(answersDTO, dto.AnswersDTO{
			ID:         int(val.ID),
			AnswerText: val.AnswerText,
			IsTrue:     val.IsTrue,
		})
	}
	return answersDTO
}

func convertQuestionToDTO(val *models.Question) *dto.FactsDTO {
	return &dto.FactsDTO{
		ID:              int(val.ID),
		QuestionContent: val.QuestionContent,
		Answers:         convertAnswerToDTO(&val.Answers),
	}
}

func convertQuestionsToDTO(facts []models.Question) []dto.FactsDTO {

	factsDTO := []dto.FactsDTO{}

	for _, val := range facts {
		dto := convertQuestionToDTO(&val)
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
func ListQuestions(c *fiber.Ctx) error {
	questions := []models.Question{}

	database.DB.Db.Preload("Answers").Find(&questions)
	factsDTO := convertQuestionsToDTO(questions)
	return c.Status(fiber.StatusOK).JSON(factsDTO)
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
func GetSingleQuestion(c *fiber.Ctx) error {
	fact := &models.Question{}
	id := c.Params("id")

	// get first fact matching the id
	result := database.DB.Db.Preload("Answers").Find(&fact, "id = ?", id)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}
	factsDTO := convertQuestionToDTO(fact)

	return c.Status(fiber.StatusOK).JSON(factsDTO)
}

// @Description Create a
// @Summary create a fact
// @Tags Facts
// @Accept json
// @Produce json
// @Param question body string true "Question"
// @Param answer body string true "Answer"
// @Success 200 {object} dto.FactsDTO
// @Failure 404 {object} string
// @Router /api/v1/facts [post]
func CreateQuestion(c *fiber.Ctx) error {
	question := new(models.Question)

	if err := c.BodyParser(question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err})
	}

	if err := business.CheckQuestion(c, question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err})
	}

	// get current user id
	user := c.Locals("user").(*dto.UserDTO)

	// set the user id
	question.UserID = user.ID

	// create the fact
	result := database.DB.Db.Create(&question)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	dto := convertQuestionToDTO(question)

	return c.Status(fiber.StatusOK).JSON(dto)

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
func DeleteQuestion(c *fiber.Ctx) error {
	questionId := c.Params("id")

	fact := &models.Question{}
	database.DB.Db.Preload("Answers").Find(fact, "id = ?", questionId)
	if fact.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "No data with that Id exists",
		})
	}

	dto := convertQuestionToDTO(fact)
	result := database.DB.Db.Select("Answers").Delete(&fact, "id = ?", questionId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No data with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(dto)

}
