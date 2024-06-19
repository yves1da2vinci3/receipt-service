package routes

import (
	"receipt-service/api/handlers"

	"github.com/gofiber/fiber/v2"
)

// ReceiptRouter is the Router for GoFiber App
func ReceiptRouter(router fiber.Router) {
	router.Get("/", handlers.GetReceipt)
}
