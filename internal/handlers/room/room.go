package roomHandler

/*
import (
	"github.com/IliyasBaish/Quizz/database"
	handlerUtils "github.com/IliyasBaish/Quizz/internal/handlers/utils"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateRoom(c *fiber.Ctx) error {
	db := database.DB
	room := new(model.Room)

	err := c.BodyParser(room)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse room data from request", "data": err})
	}

	room.Code = handlerUtils.GetCode()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot generate code", "data": err})
	}
	room.ID = uuid.New()

	err = db.Create(room).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot save room to db", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Room saved to DB", "data": room})
}*/
