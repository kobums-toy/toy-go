package router

import (
	"project/controllers/rest"
	"project/models"
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
		apiGroup.Get("/user/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.UserController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/user", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.UserController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/user", func(ctx *fiber.Ctx) error {
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
			    controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/user", func(ctx *fiber.Ctx) error {
			return ctx.SendString("Hello, put user")
		})

		apiGroup.Delete("/user", func(ctx *fiber.Ctx) error {
			return ctx.SendString("Hello, delete user")
		})
	}
}