package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"similarty-engine/handler"
	"similarty-engine/service"
)

func main() {
	app := fiber.New()

	similarityService := service.NewSimilarityService()
	simController := handler.NewRestController(&similarityService)

	app.Post("/lines/filter", simController.FilterStrings)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})
	log.Fatal(app.Listen(":80"))
}
