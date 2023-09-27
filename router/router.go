package router

import (
	questionRoutes "github.com/IliyasBaish/Quizz/internal/routes/question"
	quizzRoutes "github.com/IliyasBaish/Quizz/internal/routes/quizz"
	websocketRoutes "github.com/IliyasBaish/Quizz/internal/routes/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	questionRoutes.SetupQuestionRoutes(api)
	quizzRoutes.SetupQuizzRoutes(api)
	//teamRoutes.SetupTeamRoutes(api)
	//roomRoutes.SetupRoomRoutes(api)
	app.Static("/", "./internal/public")

	ws := app.Group("/ws")

	websocketRoutes.SetupWebsoket(ws)
}
