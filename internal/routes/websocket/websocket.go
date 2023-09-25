package websocketRoutes

import (
	websocketHandler "github.com/IliyasBaish/Quizz/internal/handlers/websocket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupWebsoket(router fiber.Router) {
	router.Get("/admin", websocket.New(websocketHandler.AdminRoom))
}
