package teamHandler

/*
import (
	"github.com/IliyasBaish/Quizz/database"
	db_interface "github.com/IliyasBaish/Quizz/internal/database"
	handlerUtils "github.com/IliyasBaish/Quizz/internal/handlers/utils"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateTeam(c *fiber.Ctx) error {
	db := database.DB
	team := new(model.Team)

	err := c.BodyParser(team)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse team from request", "data": err})
	}

	type code struct {
		Code int `json:"roomcode"`
	}

	roomcode := new(code)
	err = c.BodyParser(roomcode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse roomcode from request", "data": err})
	}

	room := db_interface.GetRoomByCode(roomcode.Code)
	if room.ID == uuid.Nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot find room with given RoomCode", "data": nil})
	}

	team.RoomID = room.ID
	team.ID = uuid.New()
	team.Code = handlerUtils.GetCode()

	err = db.Create(team).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot save team to DB", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Team saved to DB", "data": team})
}*/
