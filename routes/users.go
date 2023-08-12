package routes

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thund3rd3v/chat-app/database"
)

func Me(c *fiber.Ctx) error {
	token := string(c.Request().Header.Peek("Authorization"))

	// Trim all white spaces off token
	token = strings.TrimSpace(token)

	if token == "" {
		return c.Status(401).JSON(ErrorResponse{Message: "Unauthorized"})
	}

	user, err := database.GetUserByToken(&token)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(ErrorResponse{Message: "Unauthorized"})
		}
		return c.Status(500).JSON(ErrorResponse{Message: "Internal Server Error"})
	}

	return c.JSON(user)
}
