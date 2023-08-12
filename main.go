package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thund3rd3v/chat-app/database"
	"github.com/thund3rd3v/chat-app/routes"
)

func main() {
	// Initialize Database
	err := database.Init()
	if err != nil {
		log.Fatalln(err)
		return
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	app.Get("/status", routes.Status)

	// Auth
	app.Get("/users/me", routes.Me)
	app.Post("/auth/signup", routes.SignUp)
	app.Post("/auth/signin", routes.SignIn)

	// Start Server On Port 3000
	log.Fatal(app.Listen(":3000"))
}
