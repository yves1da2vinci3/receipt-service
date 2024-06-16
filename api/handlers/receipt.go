package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func PrintReceipt(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}

func GetReceipt(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}
