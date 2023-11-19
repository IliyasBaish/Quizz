package db_interface

import (
	"log"

	"github.com/IliyasBaish/Quizz/database"
	"github.com/IliyasBaish/Quizz/internal/model"
)

func CreateQuestion(quizzID uint, q_text string, q_answer string, time int, score int, q_group string, file string) (*model.Question, error) {
	db := database.DB
	var question model.Question

	question.QuizzID = quizzID
	question.QText = q_text
	question.QAnswer = q_answer
	question.Time = time
	question.Score = score
	question.QGroup = q_group
	question.File = file

	var quizz model.Quizz
	quizzId := question.QuizzID

	db.Find(&quizz, "id = ?", quizzId)
	log.Println(quizz.ID)
	if quizz.ID == 0 {
		return &model.Question{}, nil
	}

	err := db.Create(&question).Error
	if err != nil {
		return &model.Question{}, err
	}

	return &question, err
}

func CreateQuestions(quizzID uint, questions []model.Question) ([]*model.Question, error) {
	var questions_results = make([]*model.Question, len(questions))
	var err error
	for i := range questions {
		questions_results[i], err = CreateQuestion(quizzID, questions[i].QText, questions[i].QAnswer, questions[i].Time, questions[i].Score, questions[i].QGroup, questions[i].File)
		if err != nil {
			return questions_results, err
		}
	}
	return questions_results, nil
}
