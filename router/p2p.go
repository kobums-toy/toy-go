package router

import (
	"toysgo/controllers/p2p"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RegisterP2PRoutes(app *fiber.App) {
	// Define P2P WebSocket route
	app.Get("/p2p/ws", websocket.New(p2p.WebSocketHandler))
}
