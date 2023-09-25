package model

import (
	"github.com/gofiber/contrib/websocket"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	QuizzID uint   `json:"quizzid"`
	QText   string `json:"qtext"`
	QAnswer string `json:"qanswer"`
	Time    int    `json:"time"`
	Score   int    `json:"score"`
	QGroup  string `json:"qgroup"`
}

type Quizz struct {
	gorm.Model
	Name      string     `json:"name"`
	Questions []Question `gorm:"foreignKey:QuizzID"`
}

type Answer struct {
	QuestionId uint
	AnswerText string
	TeamCode   int
	Right      int
}

type AnswerData struct {
	Team   *websocket.Conn
	Answer *Answer
}

/*
type Room struct {
	gorm.Model
	ID      uint      `gorm:"primary key;autoIncrement"`
	QuizzID uuid.UUID `gorm:"type:uuid"`
	Quizz   Quizz     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Code    int       `gorm:"unique"`
}

type Team struct {
	gorm.Model
	ID     uint      `gorm:"primary key;autoIncrement"`
	RoomID uuid.UUID `gorm:"type:uuid" json:"roomid"`
	Room   Room      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name   string    `json:"name"`
	Code   int       `gorm:"unique" json:"code"`
	Score  int       `json:"score"`
}

type Answer struct {
	gorm.Model
	ID         uint      `gorm:"primary key;autoIncrement"`
	TeamID     uuid.UUID `gorm:"type:uuid"`
	Team       Team      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	QuestionID uuid.UUID `gorm:"type:uuid"`
	Question   Question  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AnswerText string
	Right      bool
}*/
