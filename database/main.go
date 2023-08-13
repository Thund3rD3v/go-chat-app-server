package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thund3rd3v/chat-app/structs"

	"github.com/google/uuid"
)

var database *sql.DB

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
			username TEXT,
			value TEXT NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userId, username) REFERENCES User (id, username)
		);
	`)
	}

	if err != nil {
		return err
	}

	return nil
}

func Close() {
	err := database.Close()
	if err != nil {
		log.Println(err)
	}
}

func GetMessages(amount int, offset int) ([]structs.Message, error) {
	var messages []structs.Message

	rows, err := database.Query(`
	WITH NumberedMessages AS (
		SELECT 
			m.id, m.userId, m.value, u.id AS user_id, u.username,
			ROW_NUMBER() OVER (ORDER BY m.createdAt DESC) AS row_num
		FROM Message m
		JOIN User u ON m.userId = u.id
	)
	SELECT id, userId, value, user_id, username
	FROM NumberedMessages
	WHERE row_num BETWEEN ? AND ?
	ORDER BY row_num DESC;
	`, offset+1, offset+amount)

	if err != nil {
		log.Fatalln(err)
		return messages, err
	}

	for rows.Next() {
		var message structs.Message
		var user structs.PublicUser
		err := rows.Scan(&message.Id, &message.UserId, &message.Value, &user.ID, &user.Username)
		if err != nil {
			continue // Skip this row if there's an error
		}
		message.User = user
		messages = append(messages, message)
	}

	// Do final check if there were any errors
	if err := rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func CreateMessage(userId *int, value *string) (structs.Message, error) {
	res, err := database.Exec("INSERT INTO Message (userId, value) VALUES (?, ?)", userId, value)

	if err != nil {
		return structs.Message{}, err
	}

	messageId, _ := res.LastInsertId()

	// Fetch the inserted message along with user information using a joined query
	var createdMessage structs.Message
	err = database.QueryRow(`
		SELECT
			m.id, m.userId, m.value,
			u.id, u.username
		FROM
			Message m
		JOIN
			User u ON m.userId = u.id
		WHERE
			m.id = ?`, messageId).Scan(
		&createdMessage.Id, &createdMessage.UserId, &createdMessage.Value,
		&createdMessage.User.ID, &createdMessage.User.Username)

	if err != nil {
		return structs.Message{}, err
	}

	return createdMessage, nil
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

func GetUser(username *string, password *string) (structs.User, error) {
	var user structs.User

	row := database.QueryRow("SELECT id, username, password, token FROM User WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Token)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByToken(token *string) (structs.PublicUser, error) {
	var user structs.PublicUser

	row := database.QueryRow("SELECT id, username FROM User WHERE token = ?", token)
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}
