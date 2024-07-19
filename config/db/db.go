package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AriJaya07/go-web-socket/config"
	"github.com/AriJaya07/go-web-socket/types"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	DB *sql.DB
}

func NewMySQLStorage(cfg config.Config) *MySQLStorage {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBAddress, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")

	return &MySQLStorage{DB: db}
}

// GetHistoricalMessages retrieves historical chat messages from the database
func (s *MySQLStorage) GetHistoricalMessages() ([]types.Message, error) {
	rows, err := s.DB.Query(`
		SELECT id, user_id, username, message, createdAt
		FROM chat_messages
		ORDER BY createdAt ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var (
			id        int
			userID    string
			username  string
			message   string
			createdAt time.Time
		)
		if err := rows.Scan(&id, &userID, &username, &message, &createdAt); err != nil {
			return nil, err
		}

		// Unmarshal the JSON message field
		var msg types.Message
		if err := json.Unmarshal([]byte(message), &msg); err != nil {
			return nil, err
		}
		msg.Id = id
		msg.CreatedAt = createdAt

		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// initialize the tables
	if err := s.createUserTable(); err != nil {
		return nil, err
	}
	if err := s.createTableChatMessage(); err != nil {
		return nil, err
	}

	return s.DB, nil
}

// SaveMessage stores a chat message in the database
func (s *MySQLStorage) SaveMessage(userID, username, message string, createdAt time.Time) error {
	_, err := s.DB.Exec(`
		INSERT INTO chat_messages (user_id, username, message, createdAt)
		VALUES (?, ?, ?, ?)
	`, userID, username, message, createdAt)
	return err
}

func (s *MySQLStorage) createUserTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			roleId BIGINT UNSIGNED NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createTableChatMessage() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			user_id VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}
