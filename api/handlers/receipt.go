package handlers

import (
	"fmt"
	"log"
	"receipt-service/config"
	"receipt-service/pkg/email"
	"receipt-service/pkg/entities"
	"receipt-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	eventTemplatePath      = "./pkg/receipt/templates/eventReceipt.hbs"
	logementTemplatePath   = "./pkg/receipt/templates/logementReceipt.hbs"
	defaultPDFOutputFormat = "A4"
)

func PrintReceipt(c *fiber.Ctx) error {
	QrcodeManager := new(utils.QrcodeManager)
	TemplateManager := new(utils.TemplateManager)

	receiptType := c.Query("receiptType")
	var data interface{}

	if receiptType == "logement" {
		logementData := new(entities.LogementReceiptData)
		if err := c.BodyParser(logementData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		data = logementData
		TemplateManager.SetData(data)
		QrcodeManager.SetContent(logementData.ReservationID)
	} else if receiptType == "event" {
		eventData := new(entities.EventReceiptData)
		if err := c.BodyParser(eventData); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		data = eventData
		TemplateManager.SetData(data)
		QrcodeManager.SetContent(eventData.ReservationID)
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid receipt type")
	}

	// Generate QR code
	qrCodePath, err := utils.GenerateQRCodeFile(QrcodeManager.GetContent())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate QR code")
	}
	QrcodeManager.SetPath(qrCodePath)
	// defer func() {
	// 	if err := utils.DeleteQRCodeFile(QrcodeManager.GetPath()); err != nil {
	// 		log.Printf("Failed to delete QR code file: %v", err)
	// 	}
	// }()

	// Prepare template data
	templateData := TemplateManager.GetDataMap()
	fmt.Printf("path: %s ", QrcodeManager.GetPath())
	templateData["qrcode"] = fmt.Sprintf("./%s", QrcodeManager.GetPath())

	// Determine template path
	templatePath := logementTemplatePath
	if receiptType == "event" {
		templatePath = eventTemplatePath
	}

	// Render HTML
	renderedHTML, err := utils.RenderTemplate(templatePath, templateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to render template")
	}

	// Generate PDF
	format := c.Query("format", defaultPDFOutputFormat)
	id := uuid.New()
	outputPath := fmt.Sprintf("./uploads/receipts/receipt_%s.pdf", id.String())

	err = utils.GeneratePDF(renderedHTML, format, outputPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate PDF")
	}

	fmt.Printf("PDF saved to: %s\n", outputPath)
	// Set up your Gmail service
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.GoDotEnvVariable("GMAIL_EMAIL")
	senderPassword := config.GoDotEnvVariable("GMAIL_PASSWORD")
	fmt.Printf("email := %s\n", senderEmail)
	emailService := email.NewSMTPEmailService(smtpHost, smtpPort, senderEmail, senderPassword)

	// Send an email with an attached PDF
	to := "yves.lionel.diomande@gmail.com"

	err = emailService.SendEmail(to, outputPath)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	return c.JSON(fiber.Map{"message": "PDF generated successfully"})
}

func GetReceipt(c *fiber.Ctx) error {
	path := c.Query("receiptPath")
	fmt.Printf("Receipt path: %s\n", path)
	return c.SendFile(fmt.Sprintf(".%s", path))
}
