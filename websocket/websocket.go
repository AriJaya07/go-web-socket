package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	Username string `json:"username"`
	UserID   string `json:"userId"`
}

var clients = make(map[*websocket.Conn]bool)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		var msgData Message
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("Error while unmarshalling message: %v", err)
			continue
		}

		for client := range clients {
			if err := client.WriteMessage(messageType, msg); err != nil {
				log.Printf("Error while writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
