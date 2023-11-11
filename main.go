package main

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/IliyasBaish/Quizz/config"
	"github.com/IliyasBaish/Quizz/database"
	db_interface "github.com/IliyasBaish/Quizz/internal/database"
	handlerUtils "github.com/IliyasBaish/Quizz/internal/handlers/utils"
	"github.com/IliyasBaish/Quizz/internal/model"
	ws "github.com/IliyasBaish/Quizz/internal/websocket"
	"github.com/IliyasBaish/Quizz/router"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

var rooms = make(map[int]*ws.Room) //all rooms of server

func main() {
	//engne for handling public static files
	engine := html.New("./internal/public", ".html")

	//creatig fiber server
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//middleware for allowing Cross-Origin Resource Sharing (redirecting user)
	app.Use(cors.New())

	//connect to Postgress
	database.ConnectDB()

	//serve user websocket connection
	app.Get("/ws/:roomCode", websocket.New(func(c *websocket.Conn) {
		roomCode, err := strconv.Atoi(c.Params("roomCode"))
		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte("Cannot parse roomCode"))
			c.Close()
		}
		if rooms[roomCode] == nil {
			c.WriteMessage(websocket.CloseMessage, []byte("No room with provided roomCode"))
			c.Close()
		}

		room := rooms[roomCode]

		client := ws.Client{}
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				data := strings.Split(string(message), " ")
				if data[0] == "0" {
					client = ws.Client{
						IsClosing:  false,
						Mu:         sync.Mutex{},
						Name:       data[1],
						ClientCode: handlerUtils.GetCode(),
						Conn:       c,
						Room:       room,
						AnswerChan: make(chan model.Answer),
						Score:      0,
					}
					room.ConnectClient(&client)
					break
				}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				data := strings.Split(string(message), " ")
				if data[0] == "1" {
					log.Println("answer resieved")
					questionId, _ := strconv.Atoi(data[2])
					answer := model.Answer{QuestionId: uint(questionId), AnswerText: data[1], TeamCode: client.ClientCode, Right: 0}
					client.AnswerChan <- answer
				} else if data[0] == "getteams" {
					client.Conn.WriteJSON(fiber.Map{"command": 2, "teams": room.GetClients()})
				}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	//serve admin/organizer websocket connection
	app.Get("/ws/admin/:quizzId", websocket.New(func(c *websocket.Conn) {
		quizzID, err := strconv.Atoi(c.Params("quizzId"))
		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte("Cannot parse quizzID"))
			log.Println("Cannot parse quizzID")
			c.Close()
		}
		log.Println(quizzID)
		quizz, err := db_interface.GetQuizz(uint(quizzID))
		if err != nil {
			c.WriteMessage(websocket.CloseMessage, []byte("Cannot find quizz"))
			c.Close()
		}
		log.Println(quizz.Questions[0].QText)
		code := handlerUtils.GetCode()
		for rooms[code] != nil {
			code = handlerUtils.GetCode()
		}
		room := ws.Room{
			Quizz:          quizz,
			Admin:          c,
			QuestionNum:    -1,
			ToggleQuestion: make(chan int),
			Unregister:     make(chan *websocket.Conn),
			Clients:        make(map[*websocket.Conn]*ws.Client),
			Answer:         make(chan string),
			Answers:        make(map[*ws.Client][]*model.Answer),
		}
		rooms[code] = &room
		go rooms[code].ServeRoom()
		log.Println(code)

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				data := strings.Split(string(message), " ")
				if data[0] == "start" {
					room.QuestionNum += 1
					room.ToggleQuestion <- room.QuestionNum
				} else if data[0] == "code" {
					c.WriteJSON(fiber.Map{"command": 0, "code": code})
				} else if data[0] == "score" {
					clients := room.GetClients()
					c.WriteJSON(fiber.Map{"command": 7, "teams": clients})
					room.SendScore(clients)
				} else if data[0] == "1" {
					code, _ := strconv.Atoi(data[1])
					questionID, _ := strconv.Atoi((data[3]))
					if data[2] == "1" {
						room.RateAnswer(code, true, questionID)
					} else if data[2] == "0" {
						room.RateAnswer(code, false, questionID)
					}
				} else if data[0] == "questions" {
					room.GetQuestions()
				} else if data[0] == "2" {
					id, err := strconv.Atoi(data[1])
					if err != nil {
					}
					room.SendQuestion(id)
				} else if data[0] == "3" {
					clients := room.GetClients()
					c.WriteJSON(fiber.Map{"command": 12, "teams": clients})
					room.SendFinalScore(clients)
				} else if data[0] == "5" {
					room.Close()
				}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	//setup API routes for working with qizzes
	router.SetupRoutes(app)

	//setup routes for static public files
	addr := config.Config(("ADDR"))

	//setup home page
	app.Get("/home", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("home", fiber.Map{
			"link1": "http://" + addr + "/admin",
			"link2": "http://" + addr + "/room",
			"link3": "http://" + addr + "/create",
		})
	})

	//setuop room page(for user)
	app.Get("/room", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"link1": "ws://" + addr + "/ws/",
			"link2": "http://" + addr + "/home",
		})
	})

	//setup page for creating articles
	app.Get("/create", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("createquizz", fiber.Map{
			"link1": "http://" + addr + "/api/quizz/",
			"link2": "http://" + addr + "/home",
		})
	})

	//setup room page(for admin)
	app.Get("/admin", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("createroom", fiber.Map{
			"link1": "http://" + addr + "/api/quizz/all",
			"link2": "ws://" + addr + "/ws/admin/",
			"link3": "http://" + addr + "/home",
		})
	})

	//listen connections on port 80
	app.Listen(":80")
}
