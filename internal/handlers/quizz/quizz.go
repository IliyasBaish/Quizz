package quizzHandler

import (
	"fmt"

	db_interface "github.com/IliyasBaish/Quizz/internal/database"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/fiber/v2"
)

func CreateQuizz(c *fiber.Ctx) error {
	quizz := new(model.Quizz)
	questions := []*model.Question{}

	type IncomingData struct {
		Quizz     *model.Quizz      `json:"quizz"`
		Questions *[]model.Question `json:"questions"`
	}

	data := new(IncomingData)

	fmt.Println(c.Body())

	err := c.BodyParser(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse quizz from request", "data": err})
	}

	fmt.Println(data)

	quizz, err = db_interface.CreateQuizz(data.Quizz)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot save quizz to DB", "data": err})
	}

	questions, err = db_interface.CreateQuestions(data.Quizz.ID, *data.Questions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot save qestions to DB", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "messages": "Quizz saved", "data": quizz, "questions": questions})
}

func GetQuizz(c *fiber.Ctx) error {
	type IncomingData struct {
		QuizzID uint `json:"quizzid"`
	}

	data := new(IncomingData)

	err := c.BodyParser(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse quizzid from request", "data": err})
	}

	quizz, err := db_interface.GetQuizz(data.QuizzID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot get quizz with provided id", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "messages": "Quizz saved", "data": quizz})
}

func GetQuizzs(c *fiber.Ctx) error {
	data, err := db_interface.GetQuizzs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot get quizzs from db", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "quizzes": *data})
}

func CreateQuestion(c *fiber.Ctx) error {
	return c.Render("create_question", fiber.Map{
		"link1":   "http://127.0.0.1:80/api/question/create",
		"quizzID": c.AllParams()["quizzID"],
		"group1":  c.Query("group1"),
		"group2":  c.Query("group2"),
	})
}
