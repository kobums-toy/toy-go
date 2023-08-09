package router

import (
	"project/controllers/rest"
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
			// id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
			return c.Status(200).JSON(fiber.Map{
				"code": "ok",
				"data": "hi user/${id_}",
			})
		})

		apiGroup.Get("/user", func(c *fiber.Ctx) error {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
			var controller rest.UserController
			controller.Init(c)
			controller.Index(page_, pagesize_)
			controller.Close()
			return c.JSON(controller.Result)
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