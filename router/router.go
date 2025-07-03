package router

import (
	"fmt"
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
	// 쿼리 파라미터 추출 및 로깅
	role := conn.Query("role")
	userID := conn.Query("user_id")
	userName := conn.Query("user_name")
	broadcasterID := conn.Query("broadcaster_id")

	fmt.Printf("🔗 새 WebSocket 연결 요청\n")
	fmt.Printf("  - Role: %s\n", role)
	fmt.Printf("  - User ID: %s\n", userID)
	fmt.Printf("  - User Name: %s\n", userName)
	fmt.Printf("  - Broadcaster ID: %s\n", broadcasterID)

	// 필수 파라미터 검증
	if userID == "" {
		fmt.Printf("❌ 연결 거부: user_id 누락\n")
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","data":"user_id가 필요합니다"}`))
		conn.Close()
		return
	}

	// 역할별 처리
	switch role {
	case "broadcaster":
		fmt.Printf("✅ 방송자 핸들러로 연결: %s (%s)\n", userName, userID)
		webSocketService.HandleBroadcaster(conn)

	case "viewer":
		fmt.Printf("✅ 시청자 핸들러로 연결: %s (%s)\n", userName, userID)
		webSocketService.HandleViewer(conn)

	case "viewer_list":
		fmt.Printf("✅ 목록 구독자 핸들러로 연결: %s\n", userID)
		webSocketService.HandleViewerList(conn)

	default:
		fmt.Printf("❌ 알 수 없는 역할: %s\n", role)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","data":"유효하지 않은 역할입니다"}`))
		conn.Close()
	}
}))

	apiGroup := app.Group("/api")
	// 1. 현재 방송 목록 조회
	apiGroup.Get("/broadcasts", func(c *fiber.Ctx) error {
		broadcasts := webSocketService.GetActiveBroadcasts()
		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(broadcasts),
			"data":    broadcasts,
		})
	})


	// 2. 서버 상태 조회 (전체 통계)
	apiGroup.Get("/status", func(c *fiber.Ctx) error {
		status := webSocketService.GetServerStatus()
		return c.JSON(fiber.Map{
			"success": true,
			"data":    status,
		})
	})

	// 3. 특정 방송 상세 정보
	apiGroup.Get("/broadcasts/:broadcaster_id", func(c *fiber.Ctx) error {
		broadcasterID := c.Params("broadcaster_id")
		stats := webSocketService.GetBroadcastStats(broadcasterID)
		
		if _, exists := stats["error"]; exists {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"error":   stats["error"],
			})
		}
		
		return c.JSON(fiber.Map{
			"success": true,
			"data":    stats,
		})
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
	}
}
