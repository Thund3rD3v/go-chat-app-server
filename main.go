package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thund3rd3v/chat-app/database"
	"github.com/thund3rd3v/chat-app/routes"
	"go.uber.org/zap"
)

func main() {
	// Initialize Database
	err := database.Init()
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer database.Close()

	// Create app and setup cors
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Logger
	logger, _ := zap.NewDevelopment()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// Websocket
	app.Use("/chat/ws", func(c *fiber.Ctx) error {
		// Extract the token from the request headers or query params
		token := c.Get("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		user, err := database.GetUserByToken(&token)

		if err != nil {
			// Check if user is user not found
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).JSON(routes.ErrorResponse{Message: "Unauthorized"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(routes.ErrorResponse{Message: "Internal Server Error"})
		}

		// Store validated token in c.Locals
		c.Locals("user", user)

		return c.Next()
	})
	app.Get("/chat/ws", routes.Chat)

	// Routes
	app.Get("/users/me", routes.Me)
	app.Get("/chat/messages", routes.GetMessages)
	app.Post("/auth/signup", routes.SignUp)
	app.Post("/auth/signin", routes.SignIn)

	// Start Server On Port 3000
	log.Fatal(app.Listen(":3000"))
}
