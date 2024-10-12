package router

import (
	"strconv"
	"toysgo/controllers/rest"
	"toysgo/models"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(app *fiber.App) {
	app.Get("/api/jwt", func(ctx *fiber.Ctx) error {
		email := ctx.Query("email")
		passwd := ctx.Query("passwd")
		return ctx.JSON(JwtAuth(email, passwd))
	})
	app.Get("/api/jwt/token", func(ctx *fiber.Ctx) error {
		token := ctx.Get("Authorization")
		return ctx.JSON(JwtToken(token))
	})
	apiGroup := app.Group("/api")
	apiGroup.Use(JwtAuthRequired())
	{
		apiGroup.Get("/user/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.UserController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/me", func(ctx *fiber.Ctx) error {
			token := ctx.Get("Authorization")
			return ctx.JSON(JwtMe(token))
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
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
				controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/user", func(ctx *fiber.Ctx) error {
			item_ := &models.User{}
			ctx.BodyParser(item_)
			var controller rest.UserController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
				controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/board/:id", func(ctx *fiber.Ctx) error {
			id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
			var controller rest.BoardController
			controller.Init(ctx)
			controller.Read(id_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/board", func(ctx *fiber.Ctx) error {
			page_, _ := strconv.Atoi(ctx.Query("page"))
			pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
			var controller rest.BoardController
			controller.Init(ctx)
			controller.Index(page_, pagesize_)
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Post("/board", func(ctx *fiber.Ctx) error {
			item_ := &models.Board{}
			ctx.BodyParser(item_)
			var controller rest.BoardController
			controller.Init(ctx)
			if item_ != nil {
				controller.Insert(item_)
			} else {
				controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Put("/board", func(ctx *fiber.Ctx) error {
			item_ := &models.Board{}
			ctx.BodyParser(item_)
			var controller rest.BoardController
			controller.Init(ctx)
			if item_ != nil {
				controller.Update(item_)
			} else {
				controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Delete("/board", func(ctx *fiber.Ctx) error {
			item_ := &models.Board{}
			ctx.BodyParser(item_)
			var controller rest.BoardController
			controller.Init(ctx)
			if item_ != nil {
				controller.Delete(item_)
			} else {
				controller.Result["code"] = "error"
			}
			controller.Close()
			return ctx.JSON(controller.Result)
		})
	}
}
