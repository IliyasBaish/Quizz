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

	//setup question routes
	questionRoutes.SetupQuestionRoutes(api)
	
	//setup quizz routes
	quizzRoutes.SetupQuizzRoutes(api)

	//serve public files
	app.Static("/", "./internal/public")

	//setub websocket routes
	ws := app.Group("/ws")

	websocketRoutes.SetupWebsoket(ws)
}
