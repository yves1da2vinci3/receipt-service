package utils

import (
	"os"
	"receipt-service/pkg/entities"

	"github.com/aymerick/raymond"
)

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	// Read template file
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// Parse and execute the template with data using Handlebars
	rendered, err := raymond.Render(string(tmplBytes), data)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

type TemplateManager struct {
	data interface{}
}

func (tm *TemplateManager) SetData(data interface{}) {
	tm.data = data
}

func (tm *TemplateManager) GetData() interface{} {
	return tm.data
}

func (tm *TemplateManager) GetDataMap() map[string]interface{} {
	dataMap := make(map[string]interface{})
	switch v := tm.data.(type) {
	case *entities.LogementReceiptData:
		dataMap["logementName"] = v.LogementName
		dataMap["reservationDate"] = v.ReservationDate
		dataMap["startDate"] = v.StartDate
		dataMap["endDate"] = v.EndDate
		dataMap["userName"] = v.UserName
		dataMap["reservationId"] = v.ReservationID
	case *entities.EventReceiptData:
		dataMap["eventName"] = v.EventName
		dataMap["eventDate"] = v.EventDate
		dataMap["eventLocation"] = v.EventLocation
		dataMap["userName"] = v.UserName
		dataMap["reservationId"] = v.ReservationID
	}
	return dataMap
}
