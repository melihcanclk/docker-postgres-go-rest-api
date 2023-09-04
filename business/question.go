package business

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melihcanclk/docker-postgres-go-rest-api/config"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
)

func CheckQuestion(c *fiber.Ctx, question *models.Question) error {
	// if answers's true count is not 1, return error
	trueCount := 0
	for _, val := range question.Answers {
		if val.IsTrue {
			trueCount++
		}
	}
	if trueCount != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "There must be only one true answer"})
	}

	// check answers' length
	if len(question.Answers) < config.MIN_ANSWER_COUNT {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "There must be at least two answers"})
	}

	// check answers's length
	if len(question.Answers) > config.MAX_ANSWER_COUNT {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "There must be at most six answers"})
	}

	uniqueAnswers := make(map[string]bool)
	for _, val := range question.Answers {
		if _, ok := uniqueAnswers[val.AnswerText]; ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Answers must be unique"})
		}
		uniqueAnswers[val.AnswerText] = true
	}

	return nil
}
