package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"receipt-service/pkg/entities"
	"receipt-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	eventTemplatePath      = "./pkg/receipt/templates/eventReceipt.hbs"
	logementTemplatePath   = "./pkg/receipt/templates/logementReceipt.hbs"
	defaultPDFOutputFormat = "A5"
)

func PrintReceipt(receiptType string, datum []byte) (string, error) {
	QrcodeManager := new(utils.QrcodeManager)
	TemplateManager := new(utils.TemplateManager)
	var data interface{}

	if receiptType == "logement" {
		logementData := new(entities.LogementReceiptData)
		err := json.Unmarshal(datum, &logementData)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			return "", fmt.Errorf("invalid JSON data")
		}
		data = logementData
		TemplateManager.SetData(data)
		QrcodeManager.SetContent(logementData.ReservationID)
	} else if receiptType == "event" {
		eventData := new(entities.EventReceiptData)
		err := json.Unmarshal(datum, &eventData)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			return "", fmt.Errorf("invalid JSON data")
		}
		data = eventData
		TemplateManager.SetData(data)
		QrcodeManager.SetContent(eventData.ReservationID)
	} else {
		return "", fmt.Errorf("invalid receipt type")
	}

	// Generate QR code
	qrCodePath, err := utils.GenerateQRCodeFile(QrcodeManager.GetContent())
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code")
	}
	QrcodeManager.SetPath(qrCodePath)

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
		return "", fmt.Errorf("failed to render template")
	}

	// Generate PDF
	id := uuid.New()
	outputPath := fmt.Sprintf("./uploads/receipts/receipt_%s.pdf", id.String())

	err = utils.GeneratePDF(renderedHTML, defaultPDFOutputFormat, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF")
	}

	fmt.Printf("PDF saved to: %s\n", outputPath)
	return outputPath, nil
}

func GetReceipt(c *fiber.Ctx) error {
	path := c.Query("receiptPath")
	fmt.Printf("Receipt path: %s\n", path)
	return c.SendFile(fmt.Sprintf(".%s", path))
}
