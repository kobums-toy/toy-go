package main

import (
	"log"
	"project/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	r := fiber.New()

	router.SetRouter(r)

	log.Fatal(r.Listen(":9003"))
}