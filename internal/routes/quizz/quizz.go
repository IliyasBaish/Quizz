package quizzRoutes

import (
	quizzHandler "github.com/IliyasBaish/Quizz/internal/handlers/quizz"
	"github.com/gofiber/fiber/v2"
)

func SetupQuizzRoutes(router fiber.Router) {
	quizz := router.Group("/quizz")

	quizz.Post("/", quizzHandler.CreateQuizz)
	quizz.Post("/get", quizzHandler.GetQuizz)
	quizz.Get("/all", quizzHandler.GetQuizzs)
	quizz.Get("/create/question/:quizzID", quizzHandler.CreateQuestion)
}
