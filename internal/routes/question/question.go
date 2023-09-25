package questionRoutes

import (
	questionHandler "github.com/IliyasBaish/Quizz/internal/handlers/question"
	"github.com/gofiber/fiber/v2"
)

func SetupQuestionRoutes(router fiber.Router) {
	question := router.Group("/question")

	question.Post("/", questionHandler.CreateQuestion)
	//question.Get("/", questionHandler.GetQuestions)
	//question.Get("/:questionId", questionHandler.GetQuestion)
	//question.Put("/:questionId", func(c *fiber.Ctx) error {})
	//question.Delete("/:questionId", func(c *fiber.Ctx) error {})
}
