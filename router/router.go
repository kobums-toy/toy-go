package router

import (
	"strconv"
	"toysgo/controllers/p2p"
	"toysgo/controllers/rest"
	"toysgo/models"
	"toysgo/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
	app.Get("/p2p/webrtc", websocket.New(p2p.WebSocketHandler))

	webSocketService := services.NewWebSocketService()

	app.Get("/p2p/ws", websocket.New(func(conn *websocket.Conn) {
		// WebSocket 연결 처리
		role := conn.Query("role")
		if role == "broadcaster" {
			webSocketService.SetBroadcaster(conn)
		} else if role == "viewer" {
			webSocketService.AddViewer(conn)
		} else {
			conn.Close()
		}
	}))
	apiGroup := app.Group("/api")
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
	apiGroup.Use(JwtAuthRequired())
	{
		apiGroup.Get("/oauth/token", func(ctx *fiber.Ctx) error {
			var controller rest.KakaoController
			controller.Init(ctx)
			controller.Index()
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/oauth/naver", func(ctx *fiber.Ctx) error {
			var controller rest.NaverController
			controller.Init(ctx)
			controller.Index()
			controller.Close()
			return ctx.JSON(controller.Result)
		})

		apiGroup.Get("/oauth/google", func(ctx *fiber.Ctx) error {
			var controller rest.GoogleController
			controller.Init(ctx)
			controller.Index()
			controller.Close()
			return ctx.JSON(controller.Result)
		})

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

		// apiGroup.Get("/board/:id", func(ctx *fiber.Ctx) error {
		// 	id_, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
		// 	var controller rest.BoardController
		// 	controller.Init(ctx)
		// 	controller.Read(id_)
		// 	controller.Close()
		// 	return ctx.JSON(controller.Result)
		// })

		// apiGroup.Get("/board", func(ctx *fiber.Ctx) error {
		// 	page_, _ := strconv.Atoi(ctx.Query("page"))
		// 	pagesize_, _ := strconv.Atoi(ctx.Query("pagesize"))
		// 	var controller rest.BoardController
		// 	controller.Init(ctx)
		// 	controller.Index(page_, pagesize_)
		// 	controller.Close()
		// 	return ctx.JSON(controller.Result)
		// })

		// apiGroup.Post("/board", func(ctx *fiber.Ctx) error {
		// 	item_ := &models.Board{}
		// 	ctx.BodyParser(item_)
		// 	var controller rest.BoardController
		// 	controller.Init(ctx)
		// 	if item_ != nil {
		// 		controller.Insert(item_)
		// 	} else {
		// 		controller.Result["code"] = "error"
		// 	}
		// 	controller.Close()
		// 	return ctx.JSON(controller.Result)
		// })

		// apiGroup.Put("/board", func(ctx *fiber.Ctx) error {
		// 	item_ := &models.Board{}
		// 	ctx.BodyParser(item_)
		// 	var controller rest.BoardController
		// 	controller.Init(ctx)
		// 	if item_ != nil {
		// 		controller.Update(item_)
		// 	} else {
		// 		controller.Result["code"] = "error"
		// 	}
		// 	controller.Close()
		// 	return ctx.JSON(controller.Result)
		// })

		// apiGroup.Delete("/board", func(ctx *fiber.Ctx) error {
		// 	item_ := &models.Board{}
		// 	ctx.BodyParser(item_)
		// 	var controller rest.BoardController
		// 	controller.Init(ctx)
		// 	if item_ != nil {
		// 		controller.Delete(item_)
		// 	} else {
		// 		controller.Result["code"] = "error"
		// 	}
		// 	controller.Close()
		// 	return ctx.JSON(controller.Result)
		// })
	}
}
