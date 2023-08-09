package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(r *fiber.App) {
	r.Get("/api/jwt", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	apiGroup := r.Group("/api")
	// apiGroup.Use()
	{
		apiGroup.Get("/user/:id", func(c *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
			return c.SendString(strconv.Itoa(int(id_)))
		})

		apiGroup.Get("/user", func(c *fiber.Ctx) error {
			return c.SendString("Hello, get user")
		})

		apiGroup.Post("/user", func(c *fiber.Ctx) error {
			return c.SendString("Hello, post user")
		})

		apiGroup.Put("/user", func(c *fiber.Ctx) error {
			return c.SendString("Hello, put user")
		})

		apiGroup.Delete("/user", func(c *fiber.Ctx) error {
			return c.SendString("Hello, delete user")
		})
	}
}