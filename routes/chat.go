package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/thund3rd3v/chat-app/database"
	"github.com/thund3rd3v/chat-app/structs"
)

type ActiveConnections struct {
	connections sync.Map // Use sync.Map instead of a regular map
}

var activeConnections = ActiveConnections{}

var Chat = websocket.New(func(c *websocket.Conn) {
	activeConnections.connections.Store(c, struct{}{})

	user := c.Locals("user").(structs.PublicUser)

	defer func() {
		activeConnections.connections.Delete(c)
		c.Close()
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var event ChatEvent
		if err := json.Unmarshal(message, &event); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		if event.EventType == "message" {
			createdMessage, err := database.CreateMessage(&user.ID, &event.Value)
			if err != nil {
				log.Println("Error creating message:", err)
				continue
			}

			response := ChatGlobalEvent{
				EventType: "message",
				Message:   createdMessage,
			}

			var responseJSON bytes.Buffer
			encoder := json.NewEncoder(&responseJSON)
			if err := encoder.Encode(response); err != nil {
				log.Println("JSON encoding error:", err)
				continue
			}

			activeConnections.connections.Range(func(key, value interface{}) bool {
				conn := key.(*websocket.Conn)
				go func() {
					err := conn.WriteMessage(messageType, responseJSON.Bytes())
					if err != nil {
						log.Println("Error writing message:", err)
					}
				}()
				return true
			})
		}
	}
})

func GetMessages(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	_, err := database.GetUserByToken(&token)

	if err != nil {
		// Check if user is user not found
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Message: "Unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Message: "Internal Server Error"})
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{Message: "Invalid Offset"})
	}

	messages, err := database.GetMessages(20, offset)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{Message: "Internal Server Error"})
	}

	return c.JSON(messages)
}
