package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/thund3rd3v/chat-app/database"
)

func SignUp(c *fiber.Ctx) error {
	var reqBody SignUpRequestBody

	// Parse the request body into the reqBody struct
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(ErrorResponse{Message: "Invalid request body"})
	}

	// Trim all white spaces off username & password
	username := strings.TrimSpace(reqBody.Username)
	password := strings.TrimSpace(reqBody.Password)

	// Make sure username follows standards
	if username == "" {
		return c.Status(400).JSON(ErrorResponse{Message: "Username is required"})
	}

	if len(username) < 3 {
		return c.Status(400).JSON(ErrorResponse{Message: "Username is required to be longer than or equal to 3 characters"})
	}

	if len(username) > 30 {
		return c.Status(400).JSON(ErrorResponse{Message: "Username is required to be shorter than or equal to 30 characters"})
	}

	// Make sure password follows standards
	if password == "" {
		return c.Status(400).JSON(ErrorResponse{Message: "Password is required"})
	}

	if len(password) < 8 {
		return c.Status(400).JSON(ErrorResponse{Message: "Password is required to be longer than or equal to 8 characters"})
	}

	if len(password) > 30 {
		return c.Status(400).JSON(ErrorResponse{Message: "Password is required to be shorter than or equal to 30 characters"})
	}

	_, err := database.CreateUser(&username, &password)

	if err != nil {
		// check if user with that username is already in database if there is send error
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return c.Status(409).JSON(ErrorResponse{Message: "User with that username already exists"})
			}
		}

		return c.Status(500).JSON(ErrorResponse{Message: "Internal Server Error"})
	}

	// Send success back to user
	return c.SendStatus(200)
}
