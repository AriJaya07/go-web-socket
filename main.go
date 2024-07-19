package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/AriJaya07/go-web-socket/config"
	"github.com/AriJaya07/go-web-socket/config/db"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = &sync.Mutex{}
var storage *db.MySQLStorage

// Message holds a message
type Message struct {
	MessageType int    `json:"messageType"`
	Body        string `json:"body"`
	Username    string `json:"username"`
	UserID      string `json:"userId"`
}

func reader(conn *websocket.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message := Message{MessageType: messageType, Body: string(p)}
		broadcast <- message
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		err := saveMessageToDB(msg)
		if err != nil {
			log.Println("Error saving message to DB:", err)
		}

		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(msg.MessageType, []byte(msg.Body))
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func saveMessageToDB(msg Message) error {
	_, err := storage.DB.Exec(`
		INSERT INTO chat_messages (user_id, username, message) VALUES (?, ?, ?)
	`, msg.UserID, msg.Username, msg.Body)
	return err
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	log.Println("Client Connected")
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("WebSocket Chat Server Started")

	// Initialize configuration and database
	cfg := config.InitConfig()
	storage = db.NewMySQLStorage(cfg)
	_, err := storage.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Start a goroutine to run the cleanup task every day
	go func() {
		for {
			time.Sleep(24 * time.Hour) // Wait for one day
			if err := storage.DeleteOldMessages(); err != nil {
				log.Printf("Error deleting old messages: %v", err)
			} else {
				log.Println("Old messages deleted successfully")
			}
		}
	}()

	go handleMessages()
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
