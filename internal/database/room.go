package db_interface

/*
import (
	"github.com/IliyasBaish/Quizz/database"
	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/google/uuid"
)

func GetRoomByCode(code int) model.Room {
	db := database.DB
	var room model.Room

	room.Code = code

	db.Where(&room).First(&room)

	if room.ID == uuid.Nil {
		return model.Room{ID: uuid.Nil}
	}

	return room
}

func CreateRoom(quizzId uuid.UUID)
*/
