package main

import (
	"fmt"
	"receipt-service/api/routes"
	"receipt-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New(fiber.Config{})
	routes.SetupRoutes(app)
	app.Listen(fmt.Sprintf(":%v", config.GoDotEnvVariable("PORT")))
	println("Lolo domine le monde")
}
