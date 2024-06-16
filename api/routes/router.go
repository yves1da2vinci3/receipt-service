package routes

import (
	"receipt-service/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handlers.HelloWorld)

	// receipts
	receiptRouter := api.Group("/receipts")
	ReceiptRouter(receiptRouter)
}
