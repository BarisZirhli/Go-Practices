package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func helloHandler(c *fiber.Ctx) error {
	return c.SendString("Merhaba Fiber!")
}

func Hi() {
	app := fiber.New()
	app.Get("/hello", helloHandler)
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Fiber error:", err)
	}
}
