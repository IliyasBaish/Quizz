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

var DB *gorm.DB

func ConnectDB() {
	var err error
	p := config.Config(("DB_PORT"))
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("cannot parse DB port")
	}

	dsn := fmt.Sprintf("host=%s, port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("failed to connect to DB")
	}

	fmt.Println("Connection to DB: Success")

	DB.AutoMigrate(&model.Question{})
	DB.AutoMigrate(&model.Quizz{})
	/*
		DB.AutoMigrate(&model.Room{})
		DB.AutoMigrate(&model.Team{})
		DB.AutoMigrate(&model.Answer{})
	*/
	fmt.Println("DB Migrated")
}
