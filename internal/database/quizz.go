package db_interface

import (
	"github.com/IliyasBaish/Quizz/database"
	"github.com/IliyasBaish/Quizz/internal/model"
)

func CreateQuizz(quizz *model.Quizz) (*model.Quizz, error) {
	db := database.DB

	err := db.Create(&quizz).Error
	if err != nil {
		return &model.Quizz{}, err
	}

	return quizz, nil
}

func GetQuizz(quizzID uint) (model.Quizz, error) {
	quizz := new(model.Quizz)
	err := database.DB.Model(&model.Quizz{}).Preload("Questions").Find(quizz, quizzID).Error
	if err != nil {
		return *quizz, err
	}
	return *quizz, err
}

func GetQuizzs() (*[]model.Quizz, error) {
	db := database.DB
	quizzs := new([]model.Quizz)

	result := db.Find(&quizzs)
	if result.Error != nil {
		return nil, result.Error
	}

	return quizzs, nil
}
