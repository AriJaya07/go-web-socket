package controller

import (
	"encoding/json"
	"net/http"

	"github.com/AriJaya07/go-web-socket/config/db"
	"github.com/AriJaya07/go-web-socket/types"
)

type UserController struct {
	storage *db.MySQLStorage
}

func NewUserController(storage *db.MySQLStorage) *UserController {
	return &UserController{storage: storage}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	users := []types.User{}
	rows, err := uc.storage.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
