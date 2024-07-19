package types

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Message represents a chat message
type Message struct {
	Id        int       `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"` // Adjust type as needed
}
