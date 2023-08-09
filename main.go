package main

import (
	"log"
	"project/models"
	"project/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	r := fiber.New()

	models.InitCache()

	router.SetRouter(r)

	log.Fatal(r.Listen(":9003"))
}