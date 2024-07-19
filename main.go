package main

import (
	"log"
	"net/http"

	"github.com/AriJaya07/go-web-socket/config"
	"github.com/AriJaya07/go-web-socket/config/db"
	"github.com/AriJaya07/go-web-socket/websocket"
)

func main() {
	cfg := config.InitConfig()

	// Initialize database connection
	storage := db.NewMySQLStorage(cfg)
	_, err := storage.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup WebSocket handler
	http.HandleFunc("/ws", websocket.HandleConnection)

	// Serve static files (like your index.html)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Start the server
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
