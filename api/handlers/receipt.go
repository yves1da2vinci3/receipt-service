package handlers

import (
	"fmt"
	"log"
	"receipt-service/models"
	"receipt-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PrintReceipt(c *fiber.Ctx) error {
	data := new(models.Data)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	// renderinng HTML
	renderedHTML, err := utils.RenderTemplate("./pkg/receipt/templates/template.hbs", data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	format := c.Query("format")
	fmt.Printf("format: %s", format)
	id := uuid.New()

	// Construct the file name with the UUID

	outputPath := fmt.Sprintf("./uploads/receipt_%s.pdf", id.String())

	pdfBytes, err := utils.GeneratePDF(renderedHTML, format, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}
	//
	fmt.Printf("PDF saved to: %s\n", outputPath)
	c.Set("Content-Type", "application/pdf")
	return c.Send(pdfBytes)
}

func GetReceipt(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}
