package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func exercise(a []int) []int {
	var result []int
	for _, val := range a {
		if val%2 == 0 {
			result = append(result, val)
		}
	}
	return result
}

func fiberHandler() {
	app := fiber.New()

	app.Get("/hi", func(c *fiber.Ctx) error {
		fiberResult := exercise([]int{1, 2, 3, 4, 5, 10, 15, 20, 25, 30})
		return c.SendString(fmt.Sprint(fiberResult)) // []int → string
	})

	app.Listen(":3020")
}
