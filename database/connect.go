package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/IliyasBaish/Quizz/config"
	"github.com/IliyasBaish/Quizz/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//interface for Postgress
var DB *gorm.DB

func ConnectDB() {
	var err error

	//parsing settings data from .env file
	p := config.Config(("DB_PORT"))
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("cannot parse DB port")
	}

	//connection to Postgress database
	dsn := fmt.Sprintf("host=%s, port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("failed to connect to DB")
	}

	fmt.Println("Connection to DB: Success")

	//migrate tables of dtatbase
	DB.AutoMigrate(&model.Question{})
	DB.AutoMigrate(&model.Quizz{})
	
	fmt.Println("DB Migrated")
}
