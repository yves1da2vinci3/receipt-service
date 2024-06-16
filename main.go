package main

import (
	"fmt"
	"receipt-service/api/routes"
	"receipt-service/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})
	routes.SetupRoutes(app)
	app.Listen(fmt.Sprintf(":%v", config.PORT))
	println("Lolo domine le monde")
}
