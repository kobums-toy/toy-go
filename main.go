package main

import (
	"strings"
	"time"
	"toysgo/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	log "github.com/sirupsen/logrus"
)

func main() {
	// log.SetFormatter(&log.TextFormatter{
	// 	DisableColors: false,
	// 	FullTimestamp: true,
	// })

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] | ${status} | ${latency} | ${ip}:${port} | ${method} | ${url}\n",
		TimeFormat: time.DateTime,
	}))

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Authorization, Accept",
		AllowCredentials: true,
		AllowOrigins:     "http://140.82.12.99:3000, http://localhost:3000",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			// fiber.MethodHead,
			// fiber.MethodPatch,
		}, ","),
	}))

	router.SetRouter(app)

	log.Fatal(app.Listen(":9000"))
}
