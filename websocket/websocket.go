package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AriJaya07/go-web-socket/config/db"
	"github.com/AriJaya07/go-web-socket/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var storage *db.MySQLStorage

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client Connected")
	clients[conn] = true
	defer delete(clients, conn)

	// Send historical messages to the new client
	historicalMessages, err := storage.GetHistoricalMessages()
	if err != nil {
		log.Printf("Error retrieving historical messages: %v", err)
	} else {
		for _, msg := range historicalMessages {
			messageJSON, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error marshaling message to JSON: %v", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				log.Printf("Error sending historical message: %v", err)
				conn.Close()
				return
			}
		}
	}

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		var msgData types.Message
		if err := json.Unmarshal(msg, &msgData); err != nil {
			log.Printf("Error while unmarshalling message: %v", err)
			continue
		}

		// Store the message in the database
		err = storage.SaveMessage(msgData.UserID, msgData.Username, msgData.Message, msgData.CreatedAt)
		if err != nil {
			log.Printf("Error saving message to database: %v", err)
		}

		// Broadcast the message to all connected clients
		for client := range clients {
			messageJSON, err := json.Marshal(msgData)
			if err != nil {
				log.Printf("Error marshaling message to JSON: %v", err)
				continue
			}

			if err := client.WriteMessage(messageType, messageJSON); err != nil {
				log.Printf("Error while writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
