package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Message represents a chat message
type Message struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	Body      string `json:"message"`
	CreatedAt string `json:"createdAt"` // Adjust type as needed
}
