package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/IliyasBaish/Quizz/internal/model"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Client struct {
	IsClosing  bool
	Mu         sync.Mutex
	Name       string
	ClientCode int
	Conn       *websocket.Conn
	Room       *Room
	AnswerChan chan model.Answer
	Score      int
}

type ClientList []Client

func (c ClientList) Len() int           { return len(c) }
func (c ClientList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ClientList) Less(i, j int) bool { return c[i].Score < c[j].Score }

type ClinetData struct {
	Name  string
	Code  int
	Score int
}

type RegistrationData struct {
	Connection *websocket.Conn
	Name       string
}

type Room struct {
	/*Admins 		map[*websocket.Conn]*Client
	Clients    	map[*websocket.Conn]*Client
	RegisterAdmin   	chan *websocket.Conn
	Register   	chan *websocket.Conn
	Broadcast  	chan string
	Unregister 	chan *websocket.Conn
	QuestionNum int
	Quizz 		model.Quizz
	AdminChan 	chan string*/
	Quizz          model.Quizz
	Admin          *websocket.Conn
	QuestionNum    int
	ToggleQuestion chan int
	Register       chan RegistrationData
	Unregister     chan *websocket.Conn
	Clients        map[*websocket.Conn]*Client
	Answers        map[*Client][]*model.Answer
	Answer         chan string
	AllowQuestion  bool
}

func (room *Room) ServeRoom() {
	for {
		select {
		case answer := <-room.Answer:
			log.Println("answer received:", answer)

			if err := room.Admin.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
				log.Println("write error:", err)

				room.Admin.WriteMessage(websocket.CloseMessage, []byte{})
				room.Admin.Close()
			}
		case questionID := <-room.ToggleQuestion:
			if questionID == len(room.Quizz.Questions) {
				clients, _ := json.Marshal(room.GetClients())
				room.Admin.WriteJSON(fiber.Map{"command": 6, "teams": clients})
			} else {
				log.Println("toggle question")
				log.Println(room.Clients)
				room.Admin.WriteJSON(fiber.Map{"command": 2, "questionNumber": room.QuestionNum, "question": room.Quizz.Questions[room.QuestionNum]})
				for connection, c := range room.Clients {
					go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
						c.Mu.Lock()
						defer c.Mu.Unlock()
						connection.Conn.WriteJSON(fiber.Map{"command": 3, "questionNumber": room.QuestionNum, "question": room.Quizz.Questions[room.QuestionNum]})
					}(connection, c)
				}
				room.StartQuestion()
			}
		}
	}
}

func (room *Room) SaveAnswer(client *Client, answer model.Answer) {
	fmt.Println("SaveAnswer: ", answer)
	if room.AllowQuestion {
		fmt.Println("SaveAnswer: answer allowed")
		if room.Answers[client] == nil {
			room.Answers[client] = make([]*model.Answer, len(room.Quizz.Questions))
		}
		room.Answers[client] = append(room.Answers[client], &answer)
		room.Admin.WriteJSON(fiber.Map{"command": 3, "answer": answer, "teamCode": client.ClientCode})
	} else {
		client.Conn.WriteJSON(fiber.Map{"error": "Time is out"})
	}
}

func (room *Room) RateAnswer(clinetCode int, right bool, questionID int) {
	for conn, client := range room.Clients {
		if client.ClientCode == clinetCode {
			var answer *model.Answer
			for _, val := range room.Answers[client] {
				if val != nil {
					if val.QuestionId == uint(questionID) {
						answer = val
					}
				}
			}
			if right {
				answer.Right = 1
				client.Score += room.GetQuestion(questionID).Score
				conn.WriteJSON(fiber.Map{"command": 4, "right": 1})
			} else {
				answer.Right = -1
				conn.WriteJSON(fiber.Map{"command": 4, "right": 0})
			}
			log.Println("Score", client.Score)
		}
	}
}

func (room *Room) ConnectClient(client *Client) {
	room.Clients[client.Conn] = client
	room.Admin.WriteJSON(fiber.Map{"command": 1, "teamName": client.Name, "teamCode": client.ClientCode})
	for connection, c := range room.Clients {
		go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
			c.Mu.Lock()
			defer c.Mu.Unlock()
			connection.Conn.WriteJSON(fiber.Map{"command": 1, "teamName": client.Name, "teamCode": client.ClientCode})
		}(connection, c)
	}
	go client.ServeClient()
	log.Println("connection registered")
}

func (room *Room) StartQuestion() {
	room.AllowQuestion = true
	time.Sleep(time.Duration(float64(room.Quizz.Questions[room.QuestionNum].Time) * 1e9))
	room.AllowQuestion = false
}

func (room *Room) GetClients() []ClinetData {
	data := room.Clients
	c := make(ClientList, len(data))
	i := 0
	for _, v := range data {
		c[i] = *v
		i++
	}

	sort.Sort(sort.Reverse(c))

	clients := make([]ClinetData, len(room.Clients))

	for i := range c {
		clients[i] = ClinetData{Name: c[i].Name, Code: c[i].ClientCode, Score: c[i].Score}
	}

	log.Println("clients: ", clients)
	return clients
}

func (room *Room) SendScore(clients []ClinetData) {
	for connection, c := range room.Clients {
		go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
			c.Mu.Lock()
			defer c.Mu.Unlock()
			connection.Conn.WriteJSON(fiber.Map{"command": 5, "teams": clients, "teamScore": c.Score})
		}(connection, c)
	}
}

func (room *Room) SendFinalScore(clients []ClinetData) {
	for connection, c := range room.Clients {
		go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
			c.Mu.Lock()
			defer c.Mu.Unlock()
			connection.Conn.WriteJSON(fiber.Map{"command": 7, "teams": clients, "teamScore": c.Score})
		}(connection, c)
	}
}

func (client *Client) ServeClient() {
	for {
		select {
		case answer := <-client.AnswerChan:
			client.Room.SaveAnswer(client, answer)
		}
	}
}

func (room *Room) GetQuestions() {
	room.Admin.WriteJSON(fiber.Map{"command": 10, "questions": room.Quizz.Questions})
}

func (room *Room) SendQuestion(id int) {
	question := room.GetQuestion(id)
	fmt.Println(question)
	room.Admin.WriteJSON(fiber.Map{"command": 11, "question": question})
	for connection, c := range room.Clients {
		go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
			c.Mu.Lock()
			defer c.Mu.Unlock()
			connection.Conn.WriteJSON(fiber.Map{"command": 6, "question": question})
		}(connection, c)
	}
	room.AllowQuestion = true
}

func (room *Room) GetQuestion(id int) model.Question {
	for _, question := range room.Quizz.Questions {
		if question.ID == uint(id) {
			return question
		}
	}
	return model.Question{}
}

func (room *Room) Close() {
	for connection, c := range room.Clients {
		go func(connection *websocket.Conn, c *Client) { // send to each client in parallel so we don't block on a slow client
			connection.Conn.WriteJSON(fiber.Map{"command": 8})
		}(connection, c)
	}
	room.Admin.Conn.WriteJSON(fiber.Map{"command": 13})

	room = nil
}
