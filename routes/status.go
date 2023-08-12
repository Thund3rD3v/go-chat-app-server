package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func Status(c *fiber.Ctx) error {
	return c.JSON(StatusResponse{Ok: true, Uptime: int(time.Since(startTime).Seconds())})
}
