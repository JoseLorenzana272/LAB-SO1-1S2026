package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hola, SOPES 1")
	// })

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "UP", "message": "API4 is Ready"})
	})

	log.Fatal(app.Listen(":8081"))
}
