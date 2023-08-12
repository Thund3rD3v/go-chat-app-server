package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
)

var database *sql.DB

type PublicUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func Init() error {
	var err error // Declare err here

	database, err = sql.Open("sqlite3", "./main.db")
	if err != nil {
		return err
	}

	{
		_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS User (
			id INTEGER PRIMARY KEY,
			username TEXT UNIQUE,
			password TEXT,
			token TEXT UNIQUE,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS Message (
			id INTEGER PRIMARY KEY,
			userId INTEGER NOT NULL,
			value TEXT,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId) REFERENCES User (id)
		);
	`)
	}

	if err != nil {
		return err
	}

	return nil
}

func CreateUser(username *string, password *string) (string, error) {
	token := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.MinCost)

	if err != nil {
		return token, err
	}

	_, err = database.Exec(`INSERT INTO User (username, password, token) VALUES (?, ?, ?)`, username, string(hash), token)

	if err != nil {
		return token, err
	}

	return token, nil
}

func GetToken(username *string, password *string) (string, error) {
	var token string
	var hash string

	row := database.QueryRow("SELECT token, password FROM user WHERE username = ?", username)
	err := row.Scan(&token, &hash)
	if err != nil {
		log.Fatalln(err)
		return token, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(*password))
	if err != nil {
		return token, err
	}

	return token, nil
}

func GetUserByToken(token *string) (PublicUser, error) {
	var user PublicUser

	row := database.QueryRow("SELECT id, username FROM User WHERE token = ?", token)
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}
