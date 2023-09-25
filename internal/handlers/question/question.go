package questionHandler

import (
	db_interface "github.com/IliyasBaish/Quizz/internal/database"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/fiber/v2"
)

/*
import (
	"github.com/IliyasBaish/Quizz/database"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetQuestions(c *fiber.Ctx) error {
	db := database.DB
	var questions []model.Question

	db.Find(&questions)

	if len(questions) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No questions in db", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Questions Found", "data": questions})
}

func CreateQuestion(c *fiber.Ctx) error {
	db := database.DB
	question := new(model.Question)

	err := c.BodyParser(question)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse question", "data": err})
	}

	var quizz model.Quizz
	quizzId := question.QuizzID
	db.Find(&quizz, "id = ?", quizzId)
	if quizz.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No quizz with question's quizzID", "data": nil})
	}

	question.ID = uuid.New()

	err = db.Create(&question).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot save question to DB", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Question Created", "data": question})
}

func GetQuestion(c *fiber.Ctx) error {
	db := database.DB
	var question model.Question

	id := c.Params("questionId")
	db.Find(&question, "id = ?", id)

	if question.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No question with this id", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "error", "message": "Question found", "data": question})
}*/

func CreateQuestion(c *fiber.Ctx) error {
	question := new(model.Question)

	err := c.BodyParser(question)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Cannot parse question", "data": err})
	}

	question, err = db_interface.CreateQuestion(question.QuizzID, question.QText, question.QAnswer, question.Time, question.Score, question.QGroup)

	return c.JSON(fiber.Map{"status": "success", "message": "Question Created", "data": question})
}
