package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AriJaya07/go-web-socket/config"
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

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// initialize the tables
	if err := s.createUserTable(); err != nil {
		return nil, err
	}
	if err := s.createRoleTable(); err != nil {
		return nil, err
	}
	if err := s.createHomeTable(); err != nil {
		return nil, err
	}
	if err := s.createFinanceTable(); err != nil {
		return nil, err
	}
	if err := s.createTaskTable(); err != nil {
		return nil, err
	}
	if err := s.createRentTable(); err != nil {
		return nil, err
	}
	if err := s.createMenuFoodTable(); err != nil {
		return nil, err
	}
	if err := s.createUploadFile(); err != nil {
		return nil, err
	}

	return s.DB, nil
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
			FOREIGN KEY (roleId) REFERENCES roles(id),
			UNIQUE KEY (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createRoleTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createHomeTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS home (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			userId BIGINT UNSIGNED NOT NULL,
			financeId BIGINT UNSIGNED NOT NULL,
			taskId BIGINT UNSIGNED NOT NULL,
			rentId BIGINT UNSIGNED,
			PRIMARY KEY (id), 
			FOREIGN KEY (userId) REFERENCES users(id),
			FOREIGN KEY (financeId) REFERENCES finances(id),
			FOREIGN KEY (taskId) REFERENCES tasks(id),
			FOREIGN KEY (rentId) REFERENCES rents(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createFinanceTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS finances (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			userId BIGINT UNSIGNED NOT NULL,
			name VARCHAR(500) NOT NULL,
			income VARCHAR(255) NOT NULL,
			cost VARCHAR(255) NOT NULL,
			total VARCHAR(255) NOT NULL, 
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (userId) REFERENCES users(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createTaskTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			userId BIGINT UNSIGNED NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE', 'REJECT') NOT NULL DEFAULT 'TODO',
			startDate VARCHAR(255) NOT NULL,
			endDate VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (userId) REFERENCES users(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createUploadFile() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS uploadFile (
			id BIGINT NOT NULL AUTO_INCREMENT,
			userId BIGINT UNSIGNED NOT NULL,
			filename VARCHAR(255) NOT NULL,
			filesize INT NOT NULL,
			mimeType VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (userId) REFERENCES users(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createRentTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS rents (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			roleId BIGINT UNSIGNED NOT NULL,
			name VARCHAR(255) NOT NULL,
			price VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL,
			indexNumber INT NOT NULL,
			startDate VARCHAR(255) NOT NULL,
			endDate VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (roleId) REFERENCES roles(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createMenuFoodTable() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS menuFood (
			id BIGINT NOT NULL AUTO_INCREMENT,
			roleId BIGINT UNSIGNED NOT NULL,
			name VARCHAR(255) NOT NULL,
			price VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			image VARCHAR(500) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY (roleId) REFERENCES roles(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}
