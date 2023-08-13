package routes

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thund3rd3v/chat-app/database"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c *fiber.Ctx) error {
	var reqBody SignInRequestBody

	// Parse the request body into the reqBody struct
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(ErrorResponse{Message: "Invalid request body"})
	}

	// Make sure username follows standards
	username := strings.TrimSpace(reqBody.Username)
	password := strings.TrimSpace(reqBody.Password)

	if username == "" {
		return c.Status(400).JSON(ErrorResponse{Message: "Username is required"})
	}

	// Make sure password follows standards
	if password == "" {
		return c.Status(400).JSON(ErrorResponse{Message: "Password is required"})
	}

	// Get token from database
	user, err := database.GetUser(&username, &password)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(ErrorResponse{Message: "User with that username not found make sure you have signed up"})
		}

		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.Status(401).JSON(ErrorResponse{Message: "Invalid password make sure to check your password and try again"})
		}

		return c.Status(500).JSON(ErrorResponse{Message: "Internal Server Error"})
	}

	// Send token back to user
	return c.JSON(SignInResponse{
		Id:       user.ID,
		Username: username,
		Token:    user.Token,
	})
}
