package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Message struct {
	MessageType int    `json:"messageType"`
	Body        string `json:"body"`
	Username    string `json:"username"`
}
