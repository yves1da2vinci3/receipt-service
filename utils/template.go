package utils

import (
	"os"

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
