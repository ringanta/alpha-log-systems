package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello there stranger!")
    })

	app.Post("/ssh-attempts", func(c *fiber.Ctx) error {
		log.Info(fmt.Sprintf("Receive data: %s", c.Body()))
		return c.Send(c.Body())
	})

    app.Listen(":3000")
}
